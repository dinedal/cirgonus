package query

import (
	"types"
)

func Plugin(name string, config types.CirconusConfig) interface{} {
	item, ok := config.Plugins[name]

	if ok {
		_, ok := types.Plugins[item.Type.(string)]

		if ok {
			return types.Plugins[item.Type.(string)](item.Params)
		}
	}

	return nil
}

func AllPlugins(config types.CirconusConfig) map[string]interface{} {
	retval := make(map[string]interface{})

	for key, _ := range config.Plugins {
		retval[key] = Plugin(key, config)
	}

	return retval
}
