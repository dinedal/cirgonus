package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"plugins/command"
	"plugins/cpu_usage"
	"plugins/load_average"
	"plugins/mem_usage"
	"time"
)

var plugins = map[string]func(interface{}) interface{}{
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
	Plugins []ConfigMap
}

type MeterResult struct {
	Metric string
	Type   string
	Value  interface{}
}

func load_config(config_file string) (cc CirconusConfig, err error) {
	content, err := ioutil.ReadFile(config_file)

	if err != nil {
		return cc, err
	}

	err = json.Unmarshal(content, &cc)

	for x := 0; x < len(cc.Plugins); x++ {
		item := &cc.Plugins[x]
		if item.Type == nil {
			item.Type = item.Name
		}
	}

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
			_, ok := plugins[item.Type.(string)]

			if ok {
				res, _ := json.Marshal(MeterResult{Metric: item.Name, Type: item.Type.(string), Value: plugins[item.Type.(string)](item.Params)})

				fmt.Println(string(res))
			}
		}

		time.Sleep(1 * time.Second)
	}
}
