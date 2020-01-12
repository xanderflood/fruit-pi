package device

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/xanderflood/fruit-pi-server/lib/api"
	"github.com/xanderflood/fruit-pi/lib/tools"
	"github.com/xanderflood/fruit-pi/pkg/unit"
)

func New(
	client api.Client,
	log tools.Logger,
) *Device {
	return &Device{
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

func (d *Device) Refresh(ctx context.Context) {
	if time.Since(d.lastConfigPoll) > 10*time.Second {
		d.log.Debug("refreshing list of units")
		if err := d.refreshUnits(ctx); err != nil {
			d.log.Error(err.Error())
			return
		}
	}

	d.log.Debugf("refreshing %v units", len(d.units))
	for name, unit := range d.units {
		d.log.Debugf("refreshing unit %s", name)
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

func (d *Device) refreshUnits(ctx context.Context) error {
	dState, err := d.client.GetDeviceConfig(ctx)
	if err != nil {
		d.log.Error(err)
		return fmt.Errorf("failed to get device config: %w", err)
	}

	rawConfig := map[string]unitDescription{}
	err = json.Unmarshal([]byte(*dState.Config), &rawConfig)
	if err != nil {
		return fmt.Errorf("failed to parse device config: %w", err)
	}

	d.units = map[string]unit.Unit{}
	for name, cfg := range rawConfig {
		builder, err := unit.GetBlankUnitBuilder(cfg.Type)
		if err != nil {
			return err
		}

		d.units[name], err = (*builder).BuildFromJSON([]byte(cfg.Config), d.client, d.log)
		if err != nil {
			return err
		}
		if _, ok := d.state[name]; !ok {
			d.state[name] = d.units[name].InitialState()
		}
	}

	return nil
}
