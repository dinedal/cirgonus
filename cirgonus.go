package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/syslog"
	"os"
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
		if cc.Plugins[k].Type == nil {
			old := cc.Plugins[k]
			cc.Plugins[k] = types.ConfigMap{
				Type:   k,
				Params: old.Params,
			}
		}
	}

	return cc, err
}

func checkUsage() {
	if len(os.Args) < 2 {
		fmt.Println("usage:", os.Args[0], "<config file>")
		os.Exit(2)
	}
}

func main() {
	checkUsage()

	config, err := loadConfig(os.Args[1])

	if err != nil {
		fmt.Println("Error while loading config file:", err)
		os.Exit(1)
	}

	log := initLogger(config.Facility)

	result := web.Start(config.Listen, config, log)

	log.Crit(fmt.Sprintf("Failed to start: %s", result))
}
