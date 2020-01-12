package relay

import "github.com/xanderflood/fruit-pi/pkg/gpio"

//Relay a relay module
type Relay interface {
	Set(bool)
	// On()
	// Off()
}

//Impl standard relay implementation
type Impl struct {
	pin    gpio.OutputPin
	highOn bool
}

//New control a relay
func New(pin gpio.OutputPin) *Impl {
	return &Impl{
		pin: pin,
	}
}

//On turn the relay on
func (r *Impl) Set(on bool) {
	gpio.Set(r.pin, on)
}

// //Off turn the relay off
// func (r *Impl) Off() {
// 	if r.highOn {
// 		r.pin.High()
// 	} else {
// 		r.pin.Low()
// 	}
// }
