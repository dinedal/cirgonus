# Cirgonus

Cirgonus is a go-related pun on [Circonus](http://circonus.com) and is a
metrics collector for it (and anything else that can deal with json output).

Most of the built-in collectors are linux-only for now, and probably the future
unless pull requests happen. Many plugins very likely require a 3.0 or later
kernel release due to dependence on system structs and other deep voodoo.

Cirgonus does not need to be run as root to collect its metrics.

## Querying

Currently querying is only limited to the root path and returns all metrics.
Expect this to be resolved likely before anyone but me reads this message.

e.g., if Cirgonus is running on `ubuntu.local:8000`:

```
curl http://ubuntu.local:8000
```

Will return all the metrics.

## Returned Metrics Format

An array of objects that have three elements:

* Metric -- this is the "Name" provided in the configuration file, or Type if
  Name is omitted.
* Type -- the type of metrics plugin.
* Value -- an arbitrary json-formatted value, dependent on the plugin (and in
  the case of `command`, the output of the command.)

## Config File

A config file is required to run Cirgonus. Check out `test.json` for an example
of how it should look.

### Attributes

* Listen: `host:port` (host optional) designation on where the agent should
  listen.
* Plugins: An array of plugin definitions. See `Plugins`.

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

This is a two element tuple -- the first is the number of cpus (as measured by
linux -- so hyperthreading cores are 2 cpus), and the second is a decimal value
which indicates how much of that is in use. For example:

[2, 1.5] - two cores, 150% cpu usage (1 and a half cores are in use)

### command

This is the catch-all. The command plugin runs a command, and accepts json
output from it which it builds into the result.

An example from `test.json`:

```json
    {
      "Name": "echo hi",
      "Type": "command",
      "Params": [ "echo", "{\"hi\": 1}" ]
    }
```

This results in this MeterResult:

```json
{"Metric":"echo hi","Type":"command","Value":{"hi":1}}
```

Note the value has been injected directly into the json form so that Circonus
can treat it like a proper json value.

## License

* MIT (C) Erik Hollensbe
