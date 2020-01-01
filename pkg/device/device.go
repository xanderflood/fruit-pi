package device

import (
	"context"
	"encoding/json"
	"time"

	"github.com/xanderflood/fruit-pi-server/lib/api"
	"github.com/xanderflood/fruit-pi/lib/tools"
	"github.com/xanderflood/fruit-pi/pkg/unit"
)

func New(
	client api.Client,
	log tools.Logger,
) Device {
	return Device{
		client: client,
		log:    log,
		state:  map[string]interface{}{},
	}
}

type Device struct {
	client api.Client
	log    tools.Logger

	lastConfigPoll time.Time

	units map[string]unit.Unit
	state map[string]interface{}
}

func (d Device) Refresh(ctx context.Context) {
	if time.Since(d.lastConfigPoll) > 10*time.Second {
		if ok := d.refreshUnits(ctx); !ok {
			return
		}
	}

	for name, unit := range d.units {
		state := d.state[name]

		err := unit.Refresh(state)
		if err != nil {
			d.log.Error(err)
		}
	}

}

type unitDescription struct {
	Type   string          `json:"type"`
	Config json.RawMessage `json:"config"`
}

func (d Device) refreshUnits(ctx context.Context) bool {
	dState, err := d.client.GetDeviceConfig(ctx)
	if err != nil {
		d.log.Error(err)
		return false
	}

	rawConfig := map[string]unitDescription{}
	err = json.Unmarshal([]byte(*dState.Config), &rawConfig)
	if err != nil {
		return false
	}

	d.units = map[string]unit.Unit{}
	for name, cfg := range rawConfig {
		builder := unit.GetBlankUnitBuilder(cfg.Type)

		err = json.Unmarshal([]byte(cfg.Config), &builder)
		if err != nil {
			return false
		}

		d.units[name] = (*builder).Build(d.client, d.log)
		if _, ok := d.state[name]; !ok {
			d.state[name] = d.units[name].InitialState()
		}
	}

	return true
}
