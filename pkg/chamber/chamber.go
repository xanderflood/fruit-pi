package chamber

import (
	"fmt"
	"strconv"
	"time"

	perrors "github.com/pkg/errors"

	dbsdk "github.com/xanderflood/database/lib/sdk"

	"github.com/xanderflood/fruit-pi/pkg/am2301"
	"github.com/xanderflood/fruit-pi/pkg/gpio"
	"github.com/xanderflood/fruit-pi/pkg/relay"
)

//Strategy Strategy
type Strategy interface {
	Check(am2301.State) State
	Configuration() interface{}
}

//State State
type State struct {
	Hum bool
	Fan bool
}

//Chamber a fruiting chamber module
type Chamber interface {
	Setup() error
	Refresh() (State, am2301.State, error)
}

//Impl standard chamber implementation
type Impl struct {
	name string
	db   dbsdk.SDK

	hum      relay.Relay
	fan      relay.Relay
	sensor   am2301.AM2301
	strategy Strategy

	state State
}

//New initialze `a chamber
func New(hum, fan, sensor int, strategy Strategy) *Impl {
	return &Impl{
		strategy: strategy,
		hum:      relay.New(gpio.Open(hum)),
		fan:      relay.New(gpio.Open(fan)),
		sensor:   am2301.New(gpio.Open(sensor)),
	}
}

//Refresh refresh the state error
func (c *Impl) Refresh() (State, am2301.State, error) {
	sState, err := c.sensor.Check()
	if err != nil {
		return State{}, am2301.State{}, perrors.Wrapf(err, "[chamber:%s] failed to check sensor state", c.name)
	}

	cState := c.strategy.Check(sState)
	err = c.db.Insert(c.TableName(),
		map[string]string{
			"SensorMoment": time.Now().Format(time.RFC3339),
			"RH":           strconv.FormatFloat(sState.RH, 'f', 2, 64),
			"Temp":         strconv.FormatFloat(sState.Temp, 'f', 2, 64),
			"Fan":          strconv.FormatBool(cState.Fan),
			"Hum":          strconv.FormatBool(cState.Hum),
		},
	)

	c.ensure(cState)
	return cState, sState, nil
}

func (c *Impl) ensure(state State) {
	if state.Hum {
		c.hum.On()
	} else {
		c.hum.Off()
	}

	if state.Fan {
		c.fan.On()
	} else {
		c.fan.Off()
	}
}

//Setup setup the db table
func (c *Impl) Setup() error {
	return perrors.Wrap(c.db.CreateTable(c.TableName()), "failed creating table")
}

//TableName table name
func (c *Impl) TableName() string {
	return fmt.Sprintf("fruit-pi-%s", c.name)
}
