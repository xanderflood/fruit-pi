package unit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/xanderflood/fruit-pi-server/lib/api"
	"github.com/xanderflood/fruit-pi/lib/config"
	"github.com/xanderflood/fruit-pi/lib/tools"
	"github.com/xanderflood/fruit-pi/pkg/gpio"
	"github.com/xanderflood/fruit-pi/pkg/htg3535ch"
	"github.com/xanderflood/fruit-pi/pkg/relay"
)

//SingleFanConfig is a standard unit config
type SingleFanConfig struct {
	HumidifierRelay int `json:"humidifier_relay"`
	FanRelay        int `json:"fan_rly"`

	// TemperatureCelciusADC int `json:"temp_adc"`
	// RelativeHumidityADC   int `json:"rh_adc"`

	HumOn  float64         `json:"hum_on"`
	HumOff float64         `json:"hum_off"`
	FanOn  config.Duration `json:"fan_on"`
	FanOff config.Duration `json:"fan_off"`
}

//SingleFanUnit is a standard unit implementation
type SingleFanUnit struct {
	SingleFanConfig

	htg htg3535ch.HTG3535CH
	fan relay.Relay
	hum relay.Relay

	client api.API
	log    tools.Logger
}

//SingleFanUnitState is the persistent state for the unit
type SingleFanUnitState struct {
	Humidifier     bool      `json:"humidifier"`
	Fan            bool      `json:"fan"`
	FanLastToggled time.Time `json:"fan_last_toggled"`
}

func (c SingleFanConfig) BuildFromJSON(data []byte, client api.API, log tools.Logger) (Unit, error) {
	err := json.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("failed to parse unit config: %w", err)
	}

	return SingleFanUnit{
		SingleFanConfig: c,
		client:          client,
		log:             log,
	}, nil
}

func NewSingleFanUnit(
	c SingleFanConfig,
	client api.API,
	log tools.Logger,
) SingleFanUnit {
	unit := SingleFanUnit{
		SingleFanConfig: c,
		client:          client,
		log:             log,
	}

	unit.htg = htg3535ch.HTG3535CH{}
	unit.fan = relay.New(gpio.New(c.FanRelay))
	unit.hum = relay.New(gpio.New(c.HumidifierRelay))

	return unit
}

func (c SingleFanUnit) InitialState() interface{} {
	return &SingleFanUnitState{}
}

func (c SingleFanUnit) Refresh(stateI interface{}) error {
	state := (stateI).(*SingleFanUnitState)

	sState, err := c.htg.Read()
	if err != nil {
		return fmt.Errorf("failed to check htg sensor state: %w", err)
	}

	hum, _ := sState.Humidity.Float64()
	if state.Humidifier {
		state.Humidifier = hum < c.HumOff
	} else {
		state.Humidifier = hum < c.HumOn
	}

	if state.Fan {
		if time.Since(state.FanLastToggled) > c.FanOn.Duration {
			state.Fan = false
			state.FanLastToggled = time.Now()
		}
	} else {
		if time.Since(state.FanLastToggled) > c.FanOff.Duration {
			state.Fan = true
			state.FanLastToggled = time.Now()
		}
	}

	t, _ := sState.Temperature.Float64()
	rh, _ := sState.Humidity.Float64()
	_, err = c.client.InsertReading(context.Background(), t, rh)
	if err != nil {
		return fmt.Errorf("record sensor state: %w", err)
	}

	if state.Humidifier {
		c.hum.On()
	} else {
		c.hum.Off()
	}

	if state.Fan {
		c.fan.On()
	} else {
		c.fan.Off()
	}

	return nil
}
