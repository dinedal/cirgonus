package query

import (
	"types"
)

func AllPlugins(config types.CirconusConfig) []types.MeterResult {
	retval := make([]types.MeterResult, len(config.Plugins))

	for x, item := range config.Plugins {
		_, ok := types.Plugins[item.Type.(string)]

		if ok {
			retval[x] = types.MeterResult{
				Metric: item.Name,
				Type:   item.Type.(string),
				Value:  types.Plugins[item.Type.(string)](item.Params),
			}
		}
	}

	return retval
}
