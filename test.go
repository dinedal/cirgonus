package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"plugins/load_average"
)

var plugins = map[string]func(interface{}) interface{}{
	"load_average": load_average.GetMetric,
}

type CirconusConfig struct {
	Plugins map[string]interface{}
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

	fmt.Println(os.Args[1])
	config, err := load_config(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	res, _ := json.Marshal(config)
	fmt.Println(string(res))

	for name, params := range config.Plugins {
		_, ok := plugins[name]
		if ok {
			res, _ = json.Marshal(MeterResult{Metric: name, Value: plugins[name](params)})
			fmt.Println(string(res))
		}
	}

}
