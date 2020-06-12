package device

import "github.com/xanderflood/fruit-pi/pkg/unit"

type UnitIdentifier struct {
	Kind string
	Name string
}

type DeviceState interface {
	FetchUnitState(kind, name string) (cfg unit.Config)
	FetchUnitConfig(kind, name string) (state interface{})
}

type DeviceStateStorage struct {
	config map[string]map[string]unit.Config
	state  map[string]map[string]interface{}
}
