package shell

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

func GetMetric(params interface{}) interface{} {
	array_params := params.([]interface{})
	command := make([]string, len(array_params))

	var json_out interface{}

	for i, param := range array_params {
		command[i] = param.(string)
	}

	cmd := exec.Command(command[0], command[1:]...)

	out, err := cmd.Output()

	if err != nil {
		fmt.Println("Error while gathering output for command `", command, "`:", err)
		return nil
	}

	err = json.Unmarshal(out, &json_out)

	if err != nil {
		fmt.Println("Error while marshalling content:" + string(out))
		return nil
	}

	return json_out
}
