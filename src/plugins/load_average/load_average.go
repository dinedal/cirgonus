package load_average

/*
#include <stdlib.h>
*/
import "C"

func GetMetric(params interface{}) interface{} {
	var loadavg [3]C.double

	C.getloadavg(&loadavg[0], 3)

	return loadavg
}
