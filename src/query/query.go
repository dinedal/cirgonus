package query

import (
	"types"
)

func Plugin(name string, config types.CirconusConfig) types.MeterResult {
	item, ok := config.Plugins[name]

	if ok {
		_, ok := types.Plugins[item.Type.(string)]

		if ok {
			return types.MeterResult{
				Type:  item.Type.(string),
				Value: types.Plugins[item.Type.(string)](item.Params),
			}
		}
	}

	return types.MeterResult{}
}

func AllPlugins(config types.CirconusConfig) map[string]types.MeterResult {
	retval := make(map[string]types.MeterResult)

	for key, _ := range config.Plugins {
		retval[key] = Plugin(key, config)
	}

	return retval
}
