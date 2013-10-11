package types

import (
	"logger"
	"plugins/command"
	"plugins/cpu_usage"
	"plugins/fs_usage"
	"plugins/io_usage"
	"plugins/load_average"
	"plugins/mem_usage"
	"plugins/net_usage"
	"plugins/record"
)

var Plugins = map[string]func(interface{}, *logger.Logger) interface{}{
	"load_average": load_average.GetMetric,
	"cpu_usage":    cpu_usage.GetMetric,
	"mem_usage":    mem_usage.GetMetric,
	"command":      command.GetMetric,
	"net_usage":    net_usage.GetMetric,
	"io_usage":     io_usage.GetMetric,
	"record":       record.GetMetric,
	"fs_usage":     fs_usage.GetMetric,
}

/*
How this works:

Basically, interface is expected to be nil or an array of strings which are
single parameters passed to the Params section of each json. Each element of
the array is treated as a Params line and passed straight to the generated
json. A params of "" is treated as nil because go is kind of stupid about nils.

nil means to not include the monitor. Command is a good example of a monitor we
don't want to ever try to detect.

In the load average case, our params are "", but we want to always include it.
*/

var Detectors = map[string]func() []string{
	"load_average": func() []string { return []string{} },
	"cpu_usage":    func() []string { return []string{} },
	"mem_usage":    func() []string { return []string{} },
	"command":      func() []string { return []string(nil) },
	"net_usage":    net_usage.Detect,
	"io_usage":     io_usage.Detect,
}

type ConfigMap struct {
	Type   string
	Params interface{}
}

type CirconusConfig struct {
	Listen       string
	Username     string
	Password     string
	Facility     string
	LogLevel     string
	PollInterval uint
	Plugins      map[string]ConfigMap
}
