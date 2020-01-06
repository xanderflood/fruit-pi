package htg3535ch

import (
	"fmt"
	"math"

	"github.com/xanderflood/fruit-pi/pkg/ads1115"
)

//TemperatureK represents the HTG pin for measure temperature in Kelvins
type TemperatureK struct {
	ads1115.ADS1115
	BatchResistanceOhms float64
	VCCVolts            float64
}

//NewDefaultTemperatureK creates a new TemperatureK with default wiring configuration
func NewDefaultTemperatureK(pin int) TemperatureK {
	return NewTemperatureK(pin, 10000.0, 5.0)
}

//NewTemperatureK creates a new TemperatureK with default wiring configuration
func NewTemperatureK(pin int, batchResistanceOhms, vccVolts float64) TemperatureK {
	return TemperatureK{
		ADS1115:             ads1115.New(pin),
		BatchResistanceOhms: batchResistanceOhms,
		VCCVolts:            vccVolts,
	}
}

//Read takes a reading from the underlying ADS1115 and converts the voltage
//value to a temperature reading in Kelvins.
func (s TemperatureK) Read() (float64, error) {
	v, err := s.ADS1115.ReadVoltage()
	if err != nil {
		return 0, err
	}

	ntcResistanceOhms := s.BatchResistanceOhms * v / (s.VCCVolts - v)
	fmt.Println("v_out: ", v)
	fmt.Println("resis: ", ntcResistanceOhms)
	logR := math.Log(ntcResistanceOhms)
	temp := 1 / (8.61393e-04 + 2.56377e-04*logR + 1.68055e-07*logR*logR*logR)
	return temp, nil
}

//Humidity represents the HTG pin for measure relative humidity in percent
type Humidity struct {
	ads1115.ADS1115
}

//NewHumidity creates a new Humidity
func NewHumidity(pin int) Humidity {
	return Humidity{
		ADS1115: ads1115.New(pin),
	}
}

//Read takes a reading from the underlying ADS1115 and converts the voltage
//value to a relative humidity reading in percent.
func (s Humidity) Read() (float64, error) {
	v, err := s.ADS1115.ReadVoltage()
	if err != nil {
		return 0, err
	}

	return -1.564*v*v*v + 12.05*v*v + 8.22*v - 15.6, nil
}
