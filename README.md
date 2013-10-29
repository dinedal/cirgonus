# Cirgonus

Cirgonus is a go-related pun on [Circonus](http://circonus.com) and is a
metrics collector for it (and anything else that can deal with json output). It
also comes with `cstat`, a platform-independent `iostat` alike for gathering
cirgonus metrics from many hosts.

Most of the built-in collectors are linux-only for now, and probably the future
unless pull requests happen. Many plugins very likely require a 3.0 or later
kernel release due to dependence on system structs and other deep voodoo.

Cirgonus does not need to be run as root to collect its metrics.

Unlike other collectors that use fat tools like `netstat` and `df` which can
take expensive resources on loaded systems, Cirgonus opts to use the C
interfaces directly when it can. This allows it to keep a very small footprint;
with the go runtime, it clocks in just above 5M resident and unnoticeable CPU
usage at the time of writing. The agent can sustain over 8000qps with a
benchmarking tool like `wrk`, so it will be plenty fine getting hit once per
minute, or even once per second.

## Building Cirgonus

Cirgonus due to its C dependencies must be built on a Linux box. I strongly
recommend
[godeb](http://blog.labix.org/2013/06/15/in-flight-deb-packages-of-go) for
linux users. `cstat`, however, has no such requirement, so use it on OS X or
windows if you choose.

To build, type: `make cirgonus`. To build `cstat`, type `make cstat`.

## Running Cirgonus

`cirgonus <config file|dir>` starts cirgonus with the configuration specified
-- see the file and directory sections below for more information on how the
configuration looks.

You can also use `cirgonus generate` to generate a config file based on polled
resources.

## Config File

A config file is required to run Cirgonus. Check out `test.json` for an example
of how it should look.

You can also use `cirgonus generate` to generate a configuration file from
monitors it can use and devices you have that it can monitor. This can be nice
for automated deployments.

## Config Directory

Cirgonus can also accept `conf.d` style directory configuration. There is an
example of this form in the `config_dir_example` directory in this repository.

Running cirgonus in this mode is basically the same as a configuration file,
only you point at a directory instead.

The rules are pretty basic, but surprising (sorry!):

* `main.json` is the top-level configuration.
* All other configuration files refer to plugins:
  * Metric names are named after the filename without the extension (so
    foo.json becomes metric "foo").
  * Metrics are plain JSON objects.
  * Metrics must have a "Type" and "Params" element.

### Attributes

All attributes are currently required.

* Listen: `host:port` (host optional) designation on where the agent should
  listen.
* Username: the username for basic authentication.
* Password: the password for basic authentication.
* Facility: the syslog facility to use. Controlling log levels is a function of
  your syslogd.
* Plugins: An array of plugin definitions. See `Plugins`.


## Querying

* GET querying the root will return all metrics.
* POST querying with a `{ "Name": "metric name" }` json object will return just that metric.
* PUT querying is handled by the "record" plugin. See it below.

e.g., if Cirgonus is running on `ubuntu.local:8000`:

```
curl http://ubuntu.local:8000 -d '{ "Name": "tha load" }'
```

Will return just the "tha load" metric:

```json
[0,0.01,0.05]
```

Otherwise a plain GET to the root:

```
curl http://ubuntu.local:8000
```

Will return all the metrics (Example from `test.json` configuration):

```json
{
  "cpu_usage": [0,4],
  "echo hello": {
    "hello": 1
  },
  "echo hi": {
    "hi": 1
  },
  "mem_usage": {
    "Free":845040,
    "Total":1011956,
    "Used":166916
  },
  "tha load": [0, 0.01, 0.05]
}
```

## cstat

`cstat` is a small utility to gather statistics from cirgonus agents and
display them in an `iostat`-alike manner. For example:

```
$ cstat -hosts linux-1.local,go-test.local,linux-2.local -metric "load_average"

linux-1.local: [0,0.01,0.05]
go-test.local: [0,0.01,0.05]
linux-2.local: [0,0.01,0.05]

linux-1.local: [0,0.01,0.05]
go-test.local: [0,0.01,0.05]
linux-2.local: [0,0.01,0.05]

linux-1.local: [0,0.01,0.05]
go-test.local: [0,0.01,0.05]
linux-2.local: [0,0.01,0.05]
```

## Plugins

Plugins all have a type, an optional name (type and name are equivalent if only
one or the other is supplied) and an optional set of parameters which depend on
the type of metric collected.

Plugin types are below.

### load\_average

Returns a three float tuple of the load averages in classic unix form. Takes no
parameters.

### mem\_usage

Note that this plugin uses the cached/buffers values to determine ram usage. It
also (currently) disregards swap.

Returns an object with these parameters:

* Free -- the amount of free ram.
* Total -- the total amount of ram in the machine.
* Used -- the amount of used ram.

### cpu\_usage

This is a two element tuple -- the first is a decimal value which indicates how
much of that is in use, and the second is the number of cpus (as measured by
linux -- so hyperthreading cores are 2 cpus). For example:

[1.5, 2] - two cores, 150% cpu usage (1 and a half cores are in use)

Note that due to the way jiffies are treated, this plugin imposes a minimum
granularity of 1 second.

### net\_usage

This command gathers interface statistics between the time it was last polled
for information. The parameter (A single string) is the name of the interface.

Results (hopefully self-explanatory):

```json
{
  "eth0 usage": {
    "Received (Bytes)": 2406,
    "Received (Packets)": 25,
    "Reception Errors": 0,
    "Transmission Errors": 0,
    "Transmitted (Bytes)": 1749,
    "Transmitted (Packets)": 13
  }
}
```

### io\_usage

Similar to `net_usage`, this computes the difference between polls. It provides
a variety of IO statistics and works on both disk devices and partitions,
including device mapper devices. Future patches will attempt to rein in other
devices such as tmpfs and nfs.

Results:

```json
{
  "dm-0 usage": {
    "io time (ms)": 1768,
    "iops in progress": 0,
    "reads issued": 4513,
    "reads merged": 0,
    "sectors read": 36608,
    "sectors written": 55456,
    "time reading (ms)": 1688,
    "time writing (ms)": 127156,
    "weighted io time (ms)": 128844,
    "writes completed": 5399,
    "writes merged": 0
  }
}
```

### fs\_usage

Similar to the other usage metrics, takes a mount point (such as `/`) and
returns a 4 element tuple represented as a JSON array. The first 3 elements are
unsigned integers, the last is a boolean representing the read/write status of
the filesystem.

\[ percent used, free bytes, total bytes, read/write? \]

### command

This is the catch-all. The command plugin runs a command, and accepts json
output from it which it builds into the result.

An example from `test.json`:

```json
{
  "echo hi": {
    "Type": "command",
    "Params": [ "echo", "{\"hi\": 1}" ]
  }
}
```

This results in this:

```json
{ "echo hi": {"hi":1} }
```

Note the value has been injected directly into the json form so that Circonus
can treat it like a proper json value.

### record

This allows you to inject values into cirgonus. When a PUT query is issued with
a payload that looks like:

```json
{
  "Name": "record parameter",
  "Value": { "key" : "value", "other" : ["data", "of any structure"]}
}
```

It will show up when queried with GET or POST methods just like any other
cirgonus metric.

#### Configuration

Configuration is a bit tricky; see `test.json` or see the below example. Note that
the metric name and parameter are the same. This is critical to the accurate
function of this plugin.

For a given metric named "record_example", create a file named `record_example.json`

```json
{
  "Type": "record",
  "Params": null
}
```

#### record example usage

You can try with the above configuration example with the following `curl`
commands:

To set the value:

```
curl http://cirgonus:cirgonus@localhost:8000 -X PUT -d '{ "Name": "record_example", "Value": 1 }'
```

To get the value:

```
curl http://cirgonus:cirgonus@localhost:8000 -d '{ "Name": "record_example" }'
```

Which will return `1`.

### json_poll

This allows you to poll a service at a given URL to get metrics from that
service in json format. Those metrics will be reflected to whomever polls us
for state in the given poll interval.

Parameter is just the URL. This could be used to "tunnel" metrics from
a LXC or similar embedded system by chaining cirgonus monitors.

Example:

```json
"json_poll example": {
  "Type": "json_poll",
  "Params": "http://localhost:8080"
}
```

```
$ curl http://cirgonus:cirgonus@localhost:8000 -d '{ "Name": "json_poll example" }'
{ "test": "stuff" }
# emitted by the service at http://localhost:8080
```

## License

* MIT (C) 2013 Erik Hollensbe
