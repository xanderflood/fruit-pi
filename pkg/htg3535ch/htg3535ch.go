package htg3535ch

import (
	"encoding/json"
	"os/exec"
)

//HTG3535CH is used for accessing an HTG configured on pins
//0 (temp) and 1 (hum) through the python executable
type HTG3535CH struct{}

type State struct {
	Temperature json.Number `json:"temperature"`
	Humidity    json.Number `json:"humidity"`
}

func (HTG3535CH) Read() (state State, err error) {
	cmd := exec.Command("./htg3535ch")

	// Combine stdout and stderr
	output, err := cmd.Output()
	if err != nil {
		return
	}

	err = json.Unmarshal(output, &state)
	return
}
