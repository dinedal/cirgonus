package cpu_usage

/*
#include <unistd.h>
int get_hz(void) {
  return sysconf(_SC_CLK_TCK);
}
*/
import "C"

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func get_jiffies() (jiffies int64, cpus int64) {
	content, err := ioutil.ReadFile("/proc/stat")

	if err != nil {
		fmt.Println("While processing the cpu_usage package:", err)
	}

	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		if strings.Index(line, "cpu ") == 0 {
			parts := strings.Split(line, " ")
			fmt.Println(parts)

			for x := 2; x <= 11; x++ {

				if x == 5 {
					continue
				}

				part, err := strconv.Atoi(parts[x])

				if err != nil {
					fmt.Println("Could not convert integer from string while processing cpu_usage: ", parts[x])
				}

				jiffies += int64(part)
			}
		} else if strings.Index(line, "cpu") == 0 {
			cpus++
		}
	}

	return jiffies, cpus
}

func get_jiffy_diff() (int64, int64) {
	time1, cpus := get_jiffies()
	time.Sleep(1 * time.Second)
	time2, _ := get_jiffies()

	return time2 - time1, cpus
}

func GetMetric(params interface{}) interface{} {
	diff, cpus := get_jiffy_diff()
	fmt.Println(diff, C.get_hz())
	return [2]float64{float64(cpus), (float64(diff) / float64(C.get_hz()))}
}
