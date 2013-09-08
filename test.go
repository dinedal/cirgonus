package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"plugins/cpu_usage"
	"plugins/load_average"
	"plugins/mem_usage"
	"plugins/shell"
	"time"
)

var plugins = map[string]func(interface{}) interface{}{
	"load_average": load_average.GetMetric,
	"cpu_usage":    cpu_usage.GetMetric,
	"mem_usage":    mem_usage.GetMetric,
	"shell":        shell.GetMetric,
}

type ConfigMap struct {
	Name   string
	Params interface{}
}

type CirconusConfig struct {
	Plugins []ConfigMap
}

type MeterResult struct {
	Metric string
	Value  interface{}
}

func load_config(config_file string) (cc CirconusConfig, err error) {
	content, err := ioutil.ReadFile(config_file)

	if err != nil {
		return cc, err
	}

	err = json.Unmarshal(content, &cc)

	return cc, err
}

func check_usage() {
	if len(os.Args) < 2 {
		fmt.Println("usage:", os.Args[0], "<config file>")
		os.Exit(2)
	}
}

func main() {
	check_usage()

	config, err := load_config(os.Args[1])

	if err != nil {
		fmt.Println("Error while loading config file:", err)
		os.Exit(1)
	}

	for {
		for _, item := range config.Plugins {
			_, ok := plugins[item.Name]

			if ok {
				res, _ := json.Marshal(MeterResult{Metric: item.Name, Value: plugins[item.Name](item.Params)})

				fmt.Println(string(res))
			}
		}

		time.Sleep(1 * time.Second)
	}
}
