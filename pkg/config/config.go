package config

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/xanderflood/fruit-pi/pkg/am2301"
	"github.com/xanderflood/fruit-pi/pkg/chamber"
	"github.com/xanderflood/fruit-pi/pkg/gpio"
	"github.com/xanderflood/fruit-pi/pkg/relay"
)

//Config general configuration struct
type Config struct {
	Pins     PinMapping     `json:"pins"`
	Strategy StrategyConfig `json:"strategy"`
}

//PinMapping pins
type PinMapping struct {
	Fan    int `json:"fan"`
	Hum    int `json:"hum"`
	Sensor int `json:"sensor"`
}

//Load load configuration from a stream
func Load(r io.Reader) (*Config, error) {
	buf := bytes.NewBuffer([]byte{})
	_, err := io.Copy(buf, r)
	if err != nil {
		return &Config{}, err
	}

	cfg := &Config{}
	err = json.Unmarshal(buf.Bytes(), cfg)
	return cfg, err
}

//Apply load partial configuration from stream and merge
func (cfg *Config) Apply(r io.Reader) error {
	buf := bytes.NewBuffer([]byte{})
	_, err := io.Copy(buf, r)
	if err != nil {
		return err
	}

	return json.Unmarshal(buf.Bytes(), cfg)
}

//Save write configuration to stream
func (cfg *Config) Save(w io.Writer) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(data)
	_, err = io.Copy(w, buf)
	return err
}

//Chamber build a chamber from this configuration
func (cfg *Config) Chamber() chamber.Chamber {
	return chamber.Chamber{
		Sensor:     am2301.New(gpio.New(cfg.Pins.Sensor)),
		Fan:        relay.New(gpio.New(cfg.Pins.Fan)),
		Humidifier: relay.New(gpio.New(cfg.Pins.Hum)),
	}
}
