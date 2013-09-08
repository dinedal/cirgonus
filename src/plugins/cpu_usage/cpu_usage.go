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

func getJiffies() (jiffies int64, cpus int64) {
	content, err := ioutil.ReadFile("/proc/stat")

	if err != nil {
		fmt.Println("While processing the cpu_usage package:", err)
		return 0, 0
	}

	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		if strings.Index(line, "cpu ") == 0 {
			/* cpu with no number is the aggregate of all of them -- this is what we
			 * want to parse
			 */
			parts := strings.Split(line, " ")

			/* 2 - 11 are the time aggregates */
			for x := 2; x <= 11; x++ {

				/* 5 is the idle time, which we don't want */
				if x == 5 {
					continue
				}

				/* integer all the things */
				part, err := strconv.Atoi(parts[x])

				if err != nil {
					fmt.Println("Could not convert integer from string while processing cpu_usage: ", parts[x])
					return 0, 0
				}

				jiffies += int64(part)
			}

		} else if strings.Index(line, "cpu") == 0 {
			/* cpu with a number is the specific time -- cheat and use this for the
			 * processor count since we've already read it
			 */
			cpus++
		}
	}

	return jiffies, cpus
}

func getJiffyDiff() (int64, int64) {
	time1, cpus := getJiffies()
	time.Sleep(1 * time.Second)
	time2, _ := getJiffies()

	return time2 - time1, cpus
}

func GetMetric(params interface{}) interface{} {
	diff, cpus := getJiffyDiff()
	return [2]float64{float64(cpus), (float64(diff) / float64(C.get_hz()))}
}
