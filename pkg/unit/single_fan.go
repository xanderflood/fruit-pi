package unit

import (
	"context"
	"encoding/json"
	"fmt"
	golog "log"
	"time"

	"github.com/stianeikeland/go-rpio"
	"gopkg.in/alexcesaro/statsd.v2"

	"github.com/xanderflood/fruit-pi-server/lib/api"
	"github.com/xanderflood/fruit-pi/lib/config"
	"github.com/xanderflood/fruit-pi/lib/tools"
	"github.com/xanderflood/fruit-pi/pkg/htg3535ch"
	"github.com/xanderflood/fruit-pi/pkg/relay"
)

//SingleFanConfig is a standard unit config
type SingleFanConfig struct {
	HumRelay int `json:"hum_relay"`
	FanRelay int `json:"fan_rly"`

	TemperatureCelciusADC int    `json:"temp_adc"`
	RelativeHumidityADC   int    `json:"rh_adc"`
	VoltageCalibrationADC int    `json:"vcc_adc"`
	StatsDAddress         string `json:"statsd_address"`
	Name                  string `json:"name"`

	HumOn  float64         `json:"hum_on"`
	HumOff float64         `json:"hum_off"`
	FanOn  config.Duration `json:"fan_on"`
	FanOff config.Duration `json:"fan_off"`

	InvertHumPin bool `json:"invert_hum_pin"`
	InvertFanPin bool `json:"invert_fan_pin"`
}

//SingleFanUnit is a standard unit implementation
type SingleFanUnit struct {
	SingleFanConfig
	state *SingleFanUnitState

	temp     htg3535ch.TemperatureK
	humidity htg3535ch.Humidity
	fan      relay.Relay
	hum      relay.Relay

	stats  *statsd.Client
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

	return NewSingleFanUnit(c, client, log), nil
}

func NewSingleFanUnit(
	c SingleFanConfig,
	client api.API,
	log tools.Logger,
) *SingleFanUnit {
	unit := SingleFanUnit{
		SingleFanConfig: c,
		state:           &SingleFanUnitState{},
		client:          client,
		log:             log,
		temp:            htg3535ch.NewCalibrationTemperatureK(c.TemperatureCelciusADC, c.VoltageCalibrationADC),
		humidity:        htg3535ch.NewHumidity(c.RelativeHumidityADC),
		fan:             relay.New(rpio.Pin(c.FanRelay), c.InvertFanPin),
		hum:             relay.New(rpio.Pin(c.HumRelay), c.InvertHumPin),
	}

	stats, err := statsd.New(
		statsd.Mute(len(c.StatsDAddress) == 0),
		statsd.Address(c.StatsDAddress),
		statsd.Tags(
			"name", c.Name,
		),
		statsd.ErrorHandler(func(err error) { fmt.Println(err.Error()) }),
	)
	if err != nil {
		golog.Fatal(err)
	}

	unit.stats = stats

	return &unit
}

func (c *SingleFanUnit) SetState(state interface{}) {
	c.state = state.(*SingleFanUnitState)
}
func (c *SingleFanUnit) GetState() interface{} {
	return c.state
}

func (c *SingleFanUnit) Refresh() error {
	defer c.stats.Flush()

	tempK, err := c.temp.Read()
	if err != nil {
		return fmt.Errorf("failed to check htg temperature sensor state: %w", err)
	}
	hum, err := c.humidity.Read()
	if err != nil {
		return fmt.Errorf("failed to check htg temperature sensor state: %w", err)
	}

	if c.state.Humidifier {
		c.state.Humidifier = hum < c.HumOff
	} else {
		c.state.Humidifier = hum < c.HumOn
	}

	if c.state.Fan {
		if time.Since(c.state.FanLastToggled) > c.FanOn.Duration {
			c.state.Fan = false
			c.state.FanLastToggled = time.Now()
		}
	} else {
		if time.Since(c.state.FanLastToggled) > c.FanOff.Duration {
			c.state.Fan = true
			c.state.FanLastToggled = time.Now()
		}
	}

	tempC := tempK - 273.15
	c.log.Info("Humidity    (%):", hum)
	c.log.Info("Temperature (C):", tempC)

	c.hum.Set(c.state.Humidifier)
	c.fan.Set(c.state.Fan)

	c.log.Info("hum:", c.state.Humidifier)
	c.log.Info("fan:", c.state.Fan)

	c.stats.Gauge("rh_percent", hum)
	c.stats.Gauge("temp_c", tempC)
	c.stats.Gauge("humidifier", boolToGuage(c.state.Humidifier))
	c.stats.Gauge("fan", boolToGuage(c.state.Fan))

	// TODO remove?
	_, err = c.client.InsertReading(context.Background(), tempK, hum)
	if err != nil {
		return fmt.Errorf("record sensor state: %w", err)
	}

	return nil
}

func boolToGuage(v bool) int {
	if v {
		return 1
	}
	return 0
}
