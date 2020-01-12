package unit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio"
	"github.com/xanderflood/fruit-pi-server/lib/api"
	"github.com/xanderflood/fruit-pi/lib/config"
	"github.com/xanderflood/fruit-pi/lib/tools"
	"github.com/xanderflood/fruit-pi/pkg/relay"
)

//DummyConfig is a standard unit config
type DummyConfig struct {
	HumOn  float64         `json:"hum_on"`
	HumOff float64         `json:"hum_off"`
	FanOn  config.Duration `json:"fan_on"`
	FanOff config.Duration `json:"fan_off"`

	HumidifierRelay int `json:"humidifier_relay"`
	FanRelay        int `json:"fan_rly"`

	FakeTemp float64 `json:"fake_temp"`
	FakeHum  float64 `json:"fake_hum"`
}

//DummyUnit is a standard unit implementation
type DummyUnit struct {
	DummyConfig

	fan relay.Relay
	hum relay.Relay

	client api.API
	log    tools.Logger
}

//DummyUnitState is the persistent state for the unit
type DummyUnitState struct {
	Humidifier     bool      `json:"humidifier"`
	Fan            bool      `json:"fan"`
	FanLastToggled time.Time `json:"fan_last_toggled"`
}

func (c DummyConfig) BuildFromJSON(data []byte, client api.API, log tools.Logger) (Unit, error) {
	err := json.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("failed to parse unit config: %w", err)
	}

	return NewDummyUnit(c, client, log), nil
}

func NewDummyUnit(
	c DummyConfig,
	client api.API,
	log tools.Logger,
) DummyUnit {
	unit := DummyUnit{
		DummyConfig: c,
		client:      client,
		log:         log,
	}

	unit.fan = relay.New(rpio.Pin(c.FanRelay))
	unit.hum = relay.New(rpio.Pin(c.HumidifierRelay))

	return unit
}

func (c DummyUnit) InitialState() interface{} {
	return &DummyUnitState{}
}

func (c DummyUnit) Refresh(stateI interface{}) error {
	state := (stateI).(*DummyUnitState)

	hum := c.FakeHum
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

	tK := c.FakeTemp
	_, err := c.client.InsertReading(context.Background(), tK, hum)
	if err != nil {
		return fmt.Errorf("record sensor state: %w", err)
	}

	c.log.Info("hum:", state.Humidifier)
	c.log.Info("fan:", state.Fan)

	return nil
}
