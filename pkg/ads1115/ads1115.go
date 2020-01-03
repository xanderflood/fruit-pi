package ads1115

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math"
	"os/exec"
)

type ADS1115 struct {
	pin int
}

func New(pin int) ADS1115 {
	return ADS1115{
		pin: pin,
	}
}

func (a ADS1115) ReadVoltage() (float64, error) {
	cmd := exec.Command("./htg3535ch")

	// Combine stdout and stderr
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed executing ADS1115 read: %w", err)
	}

	bs, err := base64.StdEncoding.DecodeString(string(output))
	if err != nil {
		return 0, fmt.Errorf("failed decoding ADS1115 result: %w", err)
	}

	bits := binary.LittleEndian.Uint64(bs)
	return math.Float64frombits(bits), nil
}
