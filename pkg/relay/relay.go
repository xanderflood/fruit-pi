package relay

import "github.com/xanderflood/fruit-pi/pkg/gpio"

//Relay a relay module
type Relay interface {
	On()
	Off()
}

//Impl standard relay implementation
type Impl struct {
	pin    gpio.Pin
	highOn bool
}

//NewReplay control a relay
func NewReplay(pin gpio.Pin) *Impl {
	pin.Output()

	return &Impl{
		pin: pin,
	}
}

//On turn the relay on
func (r *Impl) On() {
	if r.highOn {
		r.pin.High()
	} else {
		r.pin.Low()
	}
}

//Off turn the relay off
func (r *Impl) Off() {
	if r.highOn {
		r.pin.High()
	} else {
		r.pin.Low()
	}
}
