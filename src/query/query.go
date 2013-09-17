package query

import (
	"fmt"
	"log/syslog"
	"types"
)

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
