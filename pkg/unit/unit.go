package unit

import (
	"fmt"
	"strings"

	"github.com/xanderflood/fruit-pi-server/lib/api"
	"github.com/xanderflood/fruit-pi/lib/tools"
)

//UnitBuilder is used to encode a Unit in JSON
type UnitBuilder interface {
	BuildFromJSON(data []byte, client api.API, logger tools.Logger) (Unit, error)
}

//Unit represents some sensors and things
type Unit interface {
	Refresh(state interface{}) error
	InitialState() interface{}
}

var Units = map[string]UnitBuilder{
	"single_fan":  SingleFanConfig{},
	"dummy":       DummyConfig{},
	"relay_dummy": RelayDummyConfig{},
}

func GetBlankUnitBuilder(kind string) (*UnitBuilder, error) {
	builder, ok := Units[strings.ToLower(kind)]
	if !ok {
		return nil, fmt.Errorf("unsupported unit type identifier `%s`", kind)
	}
	return &builder, nil
}
