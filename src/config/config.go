package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"types"
)

func LoadFromFile(configFile string) (cc types.CirconusConfig, err error) {
	content, err := ioutil.ReadFile(configFile)

	if err != nil {
		return cc, err
	}

	err = json.Unmarshal(content, &cc)

	if cc.PollInterval == 0 {
		cc.PollInterval = 1
	}

	if cc.Facility == "" {
		cc.Facility = "daemon"
	}

	if cc.LogLevel == "" {
		cc.LogLevel = "info"
	}

	for k := range cc.Plugins {
		if len(cc.Plugins[k].Type) == 0 {
			old := cc.Plugins[k]
			cc.Plugins[k] = types.ConfigMap{
				Type:   k,
				Params: old.Params,
			}
		}
	}

	return cc, err
}

func Generate() {
	config := types.CirconusConfig{
		Listen:       ":8000",
		Username:     "cirgonus",
		Password:     "cirgonus",
		Facility:     "daemon",
		LogLevel:     "info",
		PollInterval: 5,
		Plugins:      make(map[string]types.ConfigMap),
	}

	for key, value := range types.Detectors {

		retval := value()

		if retval == nil {
			continue
		}

		if len(retval) == 0 {
			config.Plugins[key] = types.ConfigMap{
				Type:   key,
				Params: nil,
			}
			continue
		}

		for _, detected := range retval {
			newkey := strings.Join([]string{detected, key}, " ")
			config.Plugins[newkey] = types.ConfigMap{
				Type:   key,
				Params: detected,
			}
		}
	}

	res, err := json.MarshalIndent(config, "", "  ")

	if err != nil {
		fmt.Println("Error encountered while generating config:", err)
		os.Exit(1)
	}

	fmt.Println(string(res))
}
