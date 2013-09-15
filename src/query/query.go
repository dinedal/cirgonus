package query

import (
	"log/syslog"
	"types"
)

func Plugin(name string, config types.CirconusConfig, log *syslog.Writer) interface{} {
	item, ok := config.Plugins[name]

	if ok {
		_, ok := types.Plugins[item.Type.(string)]

		if ok {
			return types.Plugins[item.Type.(string)](item.Params, log)
		}
	}

	return nil
}

func AllPlugins(config types.CirconusConfig, log *syslog.Writer) map[string]interface{} {
	retval := make(map[string]interface{})

	for key, _ := range config.Plugins {
		retval[key] = Plugin(key, config, log)
	}

	return retval
}
