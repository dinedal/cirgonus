package json_poll

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"logger"
	"net/http"
)

func GetMetric(params interface{}, log *logger.Logger) interface{} {
	var json_out interface{}
	url := params.(string)

	resp, err := http.Get(url)

	if err != nil {
		log.Log("crit", fmt.Sprintf("Could not contact resource at URL '%s'", url))
		return nil
	}

	defer resp.Body.Close()
	out, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Log("crit", fmt.Sprintf("Error while gathering output from URL '%s'", url))
		return nil
	}

	err = json.Unmarshal(out, &json_out)

	if err != nil {
		log.Log("crit", fmt.Sprintf("Error while marshalling content: %s", string(out)))
		return nil
	}

	return json_out
}
