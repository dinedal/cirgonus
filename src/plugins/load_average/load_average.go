package load_average

/*
#include <stdlib.h>
*/
import "C"

import (
	"log/syslog"
)

func GetMetric(params interface{}, log *syslog.Writer) interface{} {
	var loadavg [3]C.double

	log.Debug("Calling getloadavg()")

	C.getloadavg(&loadavg[0], 3)

	return loadavg
}
