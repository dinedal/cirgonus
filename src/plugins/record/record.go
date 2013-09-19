package record

import (
	"log/syslog"
	"sync"
)

var recorded_metrics map[string]interface{}

var rwmutex sync.RWMutex

func GetMetric(params interface{}, log *syslog.Writer) interface{} {
	endpoint := params.(string)

	log.Debug("here")
	rwmutex.RLock()
	result := recorded_metrics[endpoint]
	rwmutex.RUnlock()

	return result
}

func RecordMetric(name string, value interface{}, log *syslog.Writer) {
	log.Debug("here")
	rwmutex.Lock()
	if recorded_metrics == nil {
		recorded_metrics = make(map[string]interface{})
	}
	recorded_metrics[name] = value
	rwmutex.Unlock()
}
