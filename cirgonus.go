package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"types"
	"web"
)

func loadConfig(configFile string) (cc types.CirconusConfig, err error) {
	content, err := ioutil.ReadFile(configFile)

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

	web.Start(config.Listen, config)
}
