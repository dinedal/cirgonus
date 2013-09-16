package types

import (
	"log/syslog"
	"plugins/command"
	"plugins/cpu_usage"
	"plugins/io_usage"
	"plugins/load_average"
	"plugins/mem_usage"
	"plugins/net_usage"
)

var Plugins = map[string]func(interface{}, *syslog.Writer) interface{}{
	"load_average": load_average.GetMetric,
	"cpu_usage":    cpu_usage.GetMetric,
	"mem_usage":    mem_usage.GetMetric,
	"command":      command.GetMetric,
	"net_usage":    net_usage.GetMetric,
	"io_usage":     io_usage.GetMetric,
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

var Detectors = map[string]func() interface{}{
	"load_average": func() interface{} { return []string{""} },
	"cpu_usage":    func() interface{} { return []string{""} },
	"mem_usage":    func() interface{} { return []string{""} },
	"command":      func() interface{} { return nil },
	"net_usage":    net_usage.Detect,
	"io_usage":     io_usage.Detect,
}

var FacilityMap = map[string]syslog.Priority{
	"kern":     syslog.LOG_KERN,
	"user":     syslog.LOG_USER,
	"mail":     syslog.LOG_MAIL,
	"daemon":   syslog.LOG_DAEMON,
	"auth":     syslog.LOG_AUTH,
	"syslog":   syslog.LOG_SYSLOG,
	"lpr":      syslog.LOG_LPR,
	"news":     syslog.LOG_NEWS,
	"uucp":     syslog.LOG_UUCP,
	"cron":     syslog.LOG_CRON,
	"authpriv": syslog.LOG_AUTHPRIV,
	"ftp":      syslog.LOG_FTP,
	"local0":   syslog.LOG_LOCAL0,
	"local1":   syslog.LOG_LOCAL1,
	"local2":   syslog.LOG_LOCAL2,
	"local3":   syslog.LOG_LOCAL3,
	"local4":   syslog.LOG_LOCAL4,
	"local5":   syslog.LOG_LOCAL5,
	"local6":   syslog.LOG_LOCAL6,
	"local7":   syslog.LOG_LOCAL7,
}

type ConfigMap struct {
	Type   interface{}
	Params interface{}
}

type CirconusConfig struct {
	Listen   string
	Username string
	Password string
	Facility string
	Plugins  map[string]ConfigMap
}
