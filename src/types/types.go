package types

import (
	"plugins/command"
	"plugins/cpu_usage"
	"plugins/load_average"
	"plugins/mem_usage"
)

var Plugins = map[string]func(interface{}) interface{}{
	"load_average": load_average.GetMetric,
	"cpu_usage":    cpu_usage.GetMetric,
	"mem_usage":    mem_usage.GetMetric,
	"command":      command.GetMetric,
}

type ConfigMap struct {
	Name   string
	Type   interface{}
	Params interface{}
}

type CirconusConfig struct {
	Listen  string
	Plugins []ConfigMap
}

type MeterResult struct {
	Metric string
	Type   string
	Value  interface{}
}
