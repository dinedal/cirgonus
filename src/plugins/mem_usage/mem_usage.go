package mem_usage

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func GetMetric(params interface{}) interface{} {
	content, err := ioutil.ReadFile("/proc/meminfo")

	var total, buffers, cached, free int

	if err != nil {
		fmt.Println("While processing the mem_usage package:", err)
		return map[string]interface{}{}
	}

	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		parts := strings.Split(line, " ")
		id := len(parts) - 2

		switch parts[0] {
		case "MemTotal:":
			total, err = strconv.Atoi(parts[id])
		case "MemFree:":
			free, err = strconv.Atoi(parts[id])
		case "Cached:":
			cached, err = strconv.Atoi(parts[id])
		case "Buffers:":
			buffers, err = strconv.Atoi(parts[id])
		}

		if err != nil {
			fmt.Println("Could not convert integer from string while processing cpu_usage: ", parts[id])
			return map[string]interface{}{}
		}
	}

	return map[string]interface{}{
		"Total": total,
		"Free":  buffers + cached + free,
		"Used":  total - (buffers + cached + free),
	}
}
