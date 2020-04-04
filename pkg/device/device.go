package device

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
	}
}

type Device struct {
	client api.Client
	log    tools.Logger

	lastConfigPoll time.Time
	deviceConfig   []byte

	units map[string]unit.Unit
}

func (d *Device) Start(
	ctx context.Context,
	file string,
	fetch time.Duration,
	refresh time.Duration,
) (err error) {
	err = d.loadDeviceConfig(ctx, file)
	if err != nil {
		return
	}

	return d.startDaemon(ctx, file, fetch, refresh)
}

func (d *Device) startDaemon(
	ctx context.Context,
	file string,
	fetch time.Duration,
	refresh time.Duration,
) error {
	err := d.tryRecoverState(ctx, file)
	if err != nil {
		return fmt.Errorf("failed loading initial state from file: %w", err)
	}

	go func() {
		fetchTicker := time.NewTicker(fetch)
		refreshTicker := time.NewTicker(refresh)

		for {
			select {
			case <-fetchTicker.C:
				d.log.Info("fetching config")
				if ok := d.tryRefreshDeviceConfig(ctx); ok {
					d.log.Info("rebuilding list of units")
					err := d.refreshUnits(ctx)
					if err != nil {
						d.log.Infof("failed refreshing units: %s", err.Error())
						continue
					}
					d.persistState(ctx, file)
				}

				d.log.Info("draining fetch")
				for len(fetchTicker.C) > 0 {
					d.log.Info("draining fetch")
					<-fetchTicker.C
				}

			case <-refreshTicker.C:
				d.log.Infof("refreshing units")
				for name, unit := range d.units {
					d.log.Infof("refreshing unit %s", name)
					err := unit.Refresh()
					if err != nil {
						d.log.Error(err)
					}
				}

				d.log.Info("draining ticker")
				for len(refreshTicker.C) > 0 {
					d.log.Info("draining ticker")
					<-refreshTicker.C
				}
			}
		}
	}()
	return nil
}

type unitDescription struct {
	Type   string          `json:"type"`
	Config json.RawMessage `json:"config"`
}

func (d *Device) tryRefreshDeviceConfig(ctx context.Context) bool {
	device, err := d.client.GetDeviceConfig(ctx)
	if err != nil {
		d.log.Error(err)
		return false
	}
	d.log.Info("fetched")
	d.deviceConfig = *device.Config
	return true
}

func (d *Device) loadDeviceConfig(ctx context.Context, file string) error {
	f, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			d.log.Infof("no config file - using empty config")
			d.deviceConfig = []byte(`{}`)
			return nil
		}
		return err
	}
	defer f.Close()

	d.deviceConfig, err = ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return nil
}

type deviceState struct {
	Config json.RawMessage        `json:"config"`
	Units  map[string]interface{} `json:"state"`
}

func (d *Device) persistState(ctx context.Context, file string) error {
	var s deviceState
	s.Config = json.RawMessage(d.deviceConfig)
	s.Units = map[string]interface{}{}

	for name, _ := range d.units {
		s.Units[name] = d.units[name].GetState()
	}

	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(s); err != nil {
		return err
	}

	if err := ioutil.WriteFile(file, buf.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

func (d *Device) tryRecoverState(ctx context.Context, file string) error {
	f, err := os.Open(file)
	if os.IsNotExist(err) {
		d.log.Infof("config file %s not found, starting from empty config", file)
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()

	var s deviceState
	err = json.NewDecoder(f).Decode(&s)
	if err != nil {
		return err
	}

	d.deviceConfig = s.Config
	err = d.refreshUnits(ctx)
	if err != nil {
		return err
	}

	for name, _ := range d.units {
		innerErr := setStateSafely(d.units[name], s.Units[name])
		if err != nil {
			err = innerErr
		}
	}
	if err != nil {
		return err
	}

	return nil
}

func (d *Device) refreshUnits(ctx context.Context) error {
	rawConfig := map[string]unitDescription{}
	err := json.Unmarshal([]byte(d.deviceConfig), &rawConfig)
	if err != nil {
		return fmt.Errorf("failed to parse device config: %w", err)
	}

	oldUnits := d.units
	newUnits := map[string]unit.Unit{}
	for name, cfg := range rawConfig {
		builder, err := unit.GetBlankUnitBuilder(cfg.Type)
		if err != nil {
			return err
		}

		unit, err := (*builder).BuildFromJSON([]byte(cfg.Config), d.client, d.log)
		if err != nil {
			return err
		}

		newUnits[name] = unit
		if oldUnit, ok := oldUnits[name]; ok {
			err = setStateSafely(unit, oldUnit.GetState())
			if err != nil {
				d.log.Error(err)
			}
		}
	}

	d.units = newUnits
	return nil
}

func setStateSafely(unit unit.Unit, state interface{}) error {
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("failed setting stored state: %v", r)
			}
		}()
		unit.SetState(state)
	}()

	return err
}
