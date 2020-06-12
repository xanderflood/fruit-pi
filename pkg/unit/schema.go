package unit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Input <-chan Value

type Output chan<- Value

type OutputSchema struct {
	NoCaching bool `json:"no_caching"`
}

type Schema struct {
	GetEmptyConfigPointer func() TypeConfig
	Inputs                map[string]struct{}
	Outputs               map[string]OutputSchema
}

type ValueIdentifier struct {
	Unit string
	Name string
}

func (c ValueIdentifier) String() string {
	return fmt.Sprintf("%s.%s", c.Unit, c.Name)
}

func (c *ValueIdentifier) UnmarshalJSON(bs []byte) error {
	parts := strings.Split(string(bs), ".")
	if len(parts) != 2 {
		return errors.New("malformed value identifier - should be <unit>.<output>")
	}

	c.Unit = parts[0]
	c.Name = parts[1]
	return nil
}

func (c ValueIdentifier) MarshalJSON() ([]byte, error) {
	return []byte(c.String()), nil
}

type TypeConfig interface {
	New() UnitV2
}

type Config struct {
	Type   string
	Config TypeConfig
	Inputs map[string]ValueIdentifier
}

func (c *Config) UnmarshalJSON(bs []byte) error {
	var configHelper struct {
		Type       string                     `json:"type"`
		ConfigJSON json.RawMessage            `json:"config"`
		Inputs     map[string]ValueIdentifier `json:"inputs"`
	}

	err := json.Unmarshal(bs, &configHelper)
	if err != nil {
		return err
	}

	schema, ok := GetSchema(configHelper.Type)
	if !ok {
		return fmt.Errorf("unregistered unit type `%s`", configHelper.Type)
	}

	if len(configHelper.ConfigJSON) == 0 {
		configHelper.ConfigJSON = []byte("{}")
	}

	c.Config = schema.GetEmptyConfigPointer()
	err = json.Unmarshal(configHelper.ConfigJSON, &c.Config)
	if err != nil {
		return fmt.Errorf("failed parsing inner unit config: %s", err)
	}

	c.Type = configHelper.Type
	c.Inputs = configHelper.Inputs
	return nil
}

// TODO should _all_ of this, including broadcast and subscription
// be moved into pkg/device?
// then pkg/unit only knows about `<-chan<- Value` and UnitV2

func (c Config) Build(
	ctx context.Context, unitName string,
	broadcastsSoFar map[string]Broadcasts,
) (unit UnitV2, newBroadcasts Broadcasts, newSubscriptions Subscriptions, err error) {
	schema, ok := GetSchema(c.Type)
	if !ok {
		err = errors.New("unregistered type TODO")
		return
	}

	newSubscriptions = Subscriptions{}
	declared := map[string]bool{}
	for inputName, valID := range c.Inputs {

		sourceUnit, ok := broadcastsSoFar[valID.Unit]
		if !ok {
			err = fmt.Errorf("reference to undeclared unit %s", valID.Unit)
			return
		}
		source, ok := sourceUnit[valID.Name]
		if !ok {
			err = fmt.Errorf("unit %s does not broadcast under the name %s", valID.Unit, valID.Name)
			return
		}

		sub := NewSubscription()
		newSubscriptions[inputName] = sub
		source.Subscribe(sub)

		declared[inputName] = true
	}

	for inputName := range schema.Inputs {
		if s, ok := declared[inputName]; !(ok && s) {
			err = errors.New("not all required inputs were declared")
			return
		}
	}

	newBroadcasts = Broadcasts{}
	for outputName, oSchema := range schema.Outputs {
		newBroadcasts[outputName] = NewBroadcast(ctx, oSchema)
	}

	unit = c.Config.New()
	return
}
