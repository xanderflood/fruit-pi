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

//TODO make the state refreshes asynchronous (but atomic)
//TODO add a timeout around the API interaction, or maybe
// just less aggressive retry rules, since it gets retried
// anyways

func New(
	client api.Client,
	log tools.Logger,
) *Device {
	return &Device{
		client: client,
		log:    log,
	}
}

type Device struct {
	client api.Client
	log    tools.Logger

	lastConfigPoll  time.Time
	lastDeviceState api.Device

	units map[string]unit.Unit
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

		//if this unit has stored state, use it.
		//this step is typically redundant.
		err := unit.Refresh()
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
	d.tryRefreshDeviceState(ctx)

	rawConfig := map[string]unitDescription{}
	err := json.Unmarshal([]byte(*d.lastDeviceState.Config), &rawConfig)
	if err != nil {
		return fmt.Errorf("failed to parse device config: %w", err)
	}

	oldUnits := d.units
	d.units = map[string]unit.Unit{}
	for name, cfg := range rawConfig {
		builder, err := unit.GetBlankUnitBuilder(cfg.Type)
		if err != nil {
			return err
		}

		unit, err := (*builder).BuildFromJSON([]byte(cfg.Config), d.client, d.log)
		if err != nil {
			return err
		}

		d.units[name] = unit
		if oldUnit, ok := oldUnits[name]; ok {
			func() {
				defer func() {
					if r := recover(); r != nil {
						err = fmt.Errorf("failed setting stored state: %v", r)
					}
					unit.SetState(oldUnit.GetState())
				}()
				return
			}()
			if err != nil {
				d.log.Error(err)
			}
		}
	}

	return nil
}

func (d *Device) tryRefreshDeviceState(ctx context.Context) {
	device, err := d.client.GetDeviceConfig(ctx)
	if err != nil {
		d.log.Error(err)
		return
	}
	d.lastDeviceState = device
}
