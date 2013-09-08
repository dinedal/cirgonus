package net_usage

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var file_pattern = "/sys/class/net/%s/statistics/"

var file_map = map[string]string{
	"rx_bytes":   "Received (Bytes)",
	"tx_bytes":   "Transmitted (Bytes)",
	"tx_errors":  "Transmission Errors",
	"rx_errors":  "Reception Errors",
	"rx_packets": "Received (Packets)",
	"tx_packets": "Transmitted (Packets)",
}

var last_metrics map[string]uint64
var rwmutex sync.RWMutex

func readFile(base_path string, metric string) (uint64, error) {
	out, err := ioutil.ReadFile(filepath.Join(base_path, metric))

	if err != nil {
		return 0, err
	}

	out_i, err := strconv.ParseUint(strings.Split(string(out), "\n")[0], 10, 64)

	return out_i, err
}

func GetMetric(params interface{}) interface{} {

	new_metrics := false

	if last_metrics == nil {
		rwmutex.Lock()
		last_metrics = make(map[string]uint64)
		rwmutex.Unlock()
		new_metrics = true
	}

	metrics := make(map[string]uint64)
	difference := make(map[string]uint64)
	device := params.(string)

	base_path := fmt.Sprintf(file_pattern, device)

	for fn, metric := range file_map {
		result, err := readFile(base_path, fn)

		if err == nil {
			metrics[metric] = result
		} else {
			metrics[metric] = 0
		}
	}

	for metric, value := range metrics {
		if new_metrics {
			difference[metric] = 0
		} else {
			rwmutex.RLock()
			difference[metric] = value - last_metrics[metric]
			rwmutex.RUnlock()
		}

		rwmutex.Lock()
		last_metrics[metric] = value
		rwmutex.Unlock()
	}

	return difference
}
