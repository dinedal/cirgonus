package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/syslog"
	"os"
	"strings"
	"types"
	"web"
)

func initLogger(facility string) *syslog.Writer {
	log, err := syslog.New(types.FacilityMap[facility]|syslog.LOG_INFO, "cirgonus")

	if err != nil {
		panic(fmt.Sprintf("Cannot connect to syslog: %s", err))
	}

	err = log.Info("Initialized Logger")
	if err != nil {
		panic(fmt.Sprintf("Cannot write to syslog: %s", err))
	}

	return log
}

func loadConfig(configFile string) (cc types.CirconusConfig, err error) {
	content, err := ioutil.ReadFile(configFile)

	if err != nil {
		return cc, err
	}

	err = json.Unmarshal(content, &cc)

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

func generateConfig() {
	config := types.CirconusConfig{
		Listen:   ":8000",
		Username: "cirgonus",
		Password: "cirgonus",
		Facility: "daemon",
		Plugins:  make(map[string]types.ConfigMap),
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

func checkUsage() {
	if len(os.Args) < 2 {
		fmt.Println("usage:", os.Args[0], "<config file> or 'generate'")
		os.Exit(2)
	}
}

func main() {
	checkUsage()

	if os.Args[1] == "generate" {
		generateConfig()
		os.Exit(0)
	} else {
		config, err := loadConfig(os.Args[1])

		if err != nil {
			fmt.Println("Error while loading config file:", err)
			os.Exit(1)
		}

		log := initLogger(config.Facility)

		result := web.Start(config.Listen, config, log)

		log.Crit(fmt.Sprintf("Failed to start: %s", result))
		log.Close()
	}
}
