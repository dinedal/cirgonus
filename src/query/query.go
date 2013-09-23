package query

import (
	"fmt"
	"log/syslog"
	"sync"
	"time"
	"types"
)

var rwmutex sync.RWMutex
var PluginResults = map[string]interface{}{}

func ResultPoller(interval int, config types.CirconusConfig, log *syslog.Writer) {
	log.Info("Starting Result Poller")

	for {
		start := time.Now()
		GetAllResults(config, log)
		duration := time.Now().Sub(start)

		if duration < time.Second*time.Duration(interval) {
			time.Sleep(duration)
		}
	}
}

func GetAllResults(config types.CirconusConfig, log *syslog.Writer) {
	rwmutex.Lock()
	PluginResults = AllPlugins(config, log)
	rwmutex.Unlock()
}

func Plugin(name string, config types.CirconusConfig, log *syslog.Writer) interface{} {
	log.Debug(fmt.Sprintf("Plugin %s Requested", name))

	item, ok := config.Plugins[name]

	if ok {
		_, ok := types.Plugins[item.Type]

		if ok {
			log.Debug(fmt.Sprintf("Plugin %s exists, running", name))
			return types.Plugins[item.Type](item.Params, log)
		}
	}

	return nil
}

func AllPlugins(config types.CirconusConfig, log *syslog.Writer) map[string]interface{} {
	retval := make(map[string]interface{})

	log.Debug("Querying All Plugins")

	for key, _ := range config.Plugins {
		retval[key] = Plugin(key, config, log)
	}

	log.Debug("Done Querying All Plugins")

	return retval
}
