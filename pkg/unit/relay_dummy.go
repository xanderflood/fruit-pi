package unit

import (
	"encoding/json"
	"fmt"

	"github.com/xanderflood/fruit-pi-server/lib/api"
	"github.com/xanderflood/fruit-pi/lib/tools"
	"github.com/xanderflood/fruit-pi/pkg/gpio"
	"github.com/xanderflood/fruit-pi/pkg/relay"
)

//RelayDummyConfig is a standard unit config
type RelayDummyConfig struct {
	HumidifierRelay int `json:"humidifier_relay"`
	FanRelay        int `json:"fan_rly"`

	HumidifierState bool `json:"humidifier_state"`
	FanState        bool `json:"fan_state"`
}

//RelayDummyUnit is a standard unit implementation
type RelayDummyUnit struct {
	RelayDummyConfig

	fan relay.Relay
	hum relay.Relay

	client api.API
	log    tools.Logger
}

func (c RelayDummyConfig) BuildFromJSON(data []byte, client api.API, log tools.Logger) (Unit, error) {
	err := json.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("failed to parse unit config: %w", err)
	}

	return NewRelayDummyUnit(c, client, log), nil
}

func NewRelayDummyUnit(
	c RelayDummyConfig,
	client api.API,
	log tools.Logger,
) RelayDummyUnit {
	unit := RelayDummyUnit{
		RelayDummyConfig: c,
		client:           client,
		log:              log,
	}

	unit.fan = relay.New(gpio.New(c.FanRelay))
	unit.hum = relay.New(gpio.New(c.HumidifierRelay))

	return unit
}

func (c RelayDummyUnit) InitialState() interface{} {
	return &struct{}{}
}

func (c RelayDummyUnit) Refresh(_ interface{}) error {
	if c.HumidifierState {
		c.hum.On()
	} else {
		c.hum.Off()
	}

	if c.FanState {
		c.fan.On()
	} else {
		c.fan.Off()
	}

	c.log.Info("hum:", c.HumidifierState)
	c.log.Info("fan:", c.FanState)

	return nil
}