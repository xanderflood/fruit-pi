package gpio

import (
	"fmt"
	"strings"

	"github.com/stianeikeland/go-rpio"
)

//State IO pin state
type State = rpio.State

//States state names
var States = map[State]string{
	Low:  "low",
	High: "high",
}

const (
	//Low signal
	Low = rpio.Low

	//High signal
	High = rpio.High
)

//ParseState parse a state from a string
func ParseState(s string) (State, error) {
	if strings.ToLower(s) == States[Low] {
		return Low, nil
	} else if strings.ToLower(s) == States[High] {
		return High, nil
	}
	return State(0), fmt.Errorf("unexpected string %s, expected HIGH or LOW", s)
}

//Setup initialize memory buffers for GPIO
func Setup() error {
	return rpio.Open()
}

//Pin minimal interface for a GPIO pin
//go:generate counterfeiter . Pin
type Pin interface {
	Set(bool)

	// High()
	// Low()

	// Input()
	// Output()
	// Read() rpio.State
}

//PinAgent wrapper around rpio.Pin
type PinAgent struct {
	rpio.Pin
}

//New open a handler for a specific GPIO pin
func New(pin int) Pin {
	return PinAgent{Pin: rpio.Pin(pin)}
}

func (pin PinAgent) Set(high bool) {
	if high {
		pin.High()
	} else {
		pin.Low()
	}
}

// //Set set the state of a GPIO pin
// func Set(pin Pin, state State) {
// 	if state == Low {
// 		pin.Low()
// 	} else {
// 		pin.High()
// 	}
// }

// //WaitChange wait until the pin reliably reads `mode`, returning the elapsed duration
// //If `timeout` ms elapse in the meantime, returns an error instead
// func WaitChange(pin Pin, mode rpio.State, timeout time.Duration) (elapsed time.Duration, ok bool) {
// 	start := time.Now()

// 	for {
// 		elapsed = time.Now().Sub(start)

// 		if elapsed > timeout {
// 			return
// 		}

// 		a := pin.Read()
// 		b := pin.Read()
// 		c := pin.Read()

// 		if (a == b) && (b == c) && (c == mode) {
// 			ok = true
// 			return
// 		}
// 	}
// }
