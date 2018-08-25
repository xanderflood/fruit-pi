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
}

//State State
type State struct {
	hum bool
	fan bool
}

//Chamber a fruiting chamber module
type Chamber interface {
	Setup() error
	Refresh() error
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
func (c *Impl) Refresh() error {
	sensorState, err := c.sensor.Check()
	if err != nil {
		return perrors.Wrapf(err, "[chamber:%s] failed to check sensor state", c.name)
	}

	err = c.db.Insert(c.TableName(),
		map[string]string{
			"SensorMoment": time.Now().Format(time.RFC3339),
			"RH":           strconv.FormatFloat(sensorState.RH, 'f', 2, 64),
			"Temp":         strconv.FormatFloat(sensorState.Temp, 'f', 2, 64),
		},
	)

	c.ensure(c.strategy.Check(sensorState))
	return nil
}

func (c *Impl) ensure(state State) {
	if state.hum {
		c.hum.On()
	} else {
		c.hum.Off()
	}

	if state.fan {
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
