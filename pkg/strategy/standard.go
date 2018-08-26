package strategy

import (
	"time"

	"github.com/xanderflood/fruit-pi/lib/config"
	"github.com/xanderflood/fruit-pi/pkg/am2301"
	"github.com/xanderflood/fruit-pi/pkg/chamber"
)

type Standard struct {
	cur      chamber.State
	fanTimer *time.Timer

	Config struct {
		HumOn  float64         `json:"hum_on"`
		HumOff float64         `json:"hum_off"`
		FanOn  config.Duration `json:"fan_on"`
		FanOff config.Duration `json:"fan_off"`
	}
}

func (s *Standard) Configuration() interface{} {
	return &s.Config
}

func (s *Standard) Check(sState am2301.State) chamber.State {
	// set humidifier state
	if s.cur.Hum {
		if sState.RH > s.Config.HumOn {
			s.cur.Hum = true
		}
	} else {
		if sState.RH < s.Config.HumOff {
			s.cur.Hum = false
		}
	}

	// initialize the timer if necessary
	if s.fanTimer == nil {
		// start with an off period
		s.fanTimer = time.NewTimer(s.Config.FanOff.Duration)
	}

	// if the timer has finished, switch the fan state
	select {
	case <-s.fanTimer.C:
		if s.cur.Fan {
			s.fanTimer = time.NewTimer(s.Config.FanOff.Duration)
		} else {
			s.fanTimer = time.NewTimer(s.Config.FanOn.Duration)
		}

		s.cur.Fan = !s.cur.Fan
	default:
		// no change yet
	}

	return s.cur
}
