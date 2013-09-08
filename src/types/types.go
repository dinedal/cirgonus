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
	Type   interface{}
	Params interface{}
}

type CirconusConfig struct {
	Listen  string
	Plugins map[string]ConfigMap
}

type MeterResult struct {
	Type  string
	Value interface{}
}
