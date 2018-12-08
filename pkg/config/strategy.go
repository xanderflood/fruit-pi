package config

import (
	"encoding/json"
	"fmt"

	perrors "github.com/pkg/errors"

	"github.com/xanderflood/fruit-pi/pkg/chamber"
)

//StrategyConfigContents StrategyConfigContents
type StrategyConfigContents struct {
	Name   string            `json:"name"`
	Config chamber.Validator `json:"config"`
}

//StrategyConfig StrategyConfig
type StrategyConfig struct {
	*StrategyConfigContents

	Object chamber.Strategy
}

//UnmarshalJSON UnmarshalJSON
func (s *StrategyConfig) UnmarshalJSON(b []byte) error {
	if s.StrategyConfigContents == nil {
		s.StrategyConfigContents = &StrategyConfigContents{}
	}

	contents := s.StrategyConfigContents

	//unmarshals just the strategy name, since contents.Config is nil
	err := json.Unmarshal(b, contents)
	if err != nil {
		return err
	}

	//get a strategy object and its config object
	s.Object = chamber.Load(contents.Name)
	if s.Object == nil {
		return fmt.Errorf("invalid strategy name: %s", contents.Name)
	}
	contents.Config = s.Object.Configuration()
	err = contents.Config.Validate()
	if err != nil {
		return perrors.Wrapf(err, "invalid config for strategy type `%s`", contents.Name)
	}

	//unmarshal again, this time with a non-nil contents.Config
	err = json.Unmarshal(b, contents)
	if err != nil {
		return err
	}

	return nil
}
