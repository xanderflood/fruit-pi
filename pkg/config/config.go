package config

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/xanderflood/fruit-pi/lib/config"
)

type Config struct {
	FanPin    int             `json:"fan_pin"`
	HumPin    int             `json:"hum_pin"`
	SensorPin int             `json:"sensor_pin"`
	Interval  config.Duration `json:"interval"`
	Strategy  StrategyConfig  `json:"strategy"`
}

func Load(r io.Reader) (Config, error) {
	buf := bytes.NewBuffer([]byte{})
	_, err := io.Copy(buf, r)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{}
	err = json.Unmarshal(buf.Bytes(), &cfg)
	return cfg, err
}

func (cfg *Config) Apply(r io.Reader) error {
	buf := bytes.NewBuffer([]byte{})
	_, err := io.Copy(buf, r)
	if err != nil {
		return err
	}

	return json.Unmarshal(buf.Bytes(), cfg)
}

func (cfg *Config) Save(w io.Writer) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(data)
	_, err = io.Copy(w, buf)
	return err
}
