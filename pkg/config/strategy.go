package config

import (
	"encoding/json"
	"fmt"

	"github.com/xanderflood/fruit-pi/pkg/chamber"
	"github.com/xanderflood/fruit-pi/pkg/strategy"
)

type StrategyConfigContents struct {
	Name   string      `json:"name"`
	Config interface{} `json:"config"`
}

type StrategyConfig struct {
	*StrategyConfigContents

	Object chamber.Strategy
}

func (s *StrategyConfig) UnmarshalJSON(b []byte) error {
	if s.StrategyConfigContents == nil {
		s.StrategyConfigContents = &StrategyConfigContents{}
	}

	contents := s.StrategyConfigContents

	//unmarshals just the strategy name, since
	//contents.Config is nil
	err := json.Unmarshal(b, contents)
	if err != nil {
		return err
	}

	//get a strategy object and its config object
	s.Object = loadStrategy(contents.Name)
	if s.Object == nil {
		return fmt.Errorf("invalid strategy name: %s", contents.Name)
	}
	contents.Config = s.Object.Configuration()

	//unmarshal again, this time with a non-nil
	//contents.Config
	err = json.Unmarshal(b, contents)
	if err != nil {
		return err
	}

	return nil
}

func loadStrategy(name string) chamber.Strategy {
	switch name {
	case "standard":
		return &strategy.Standard{}
	default:
		return nil
	}
}
