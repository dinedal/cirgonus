package main

import (
	"encoding/json"
	"fmt"
	"plugins/load_average"
)

type MeterResult struct {
	Metric string
	Value  interface{}
}

func main() {
	res, _ := json.Marshal(MeterResult{Metric: "load average", Value: load_average.GetMetric()})
	fmt.Println(string(res))
}
