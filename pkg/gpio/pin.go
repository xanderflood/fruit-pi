package gpio

import (
	"errors"
	"time"

	"github.com/stianeikeland/go-rpio"
)

const (
	//Low signal
	Low = rpio.Low

	//High signal
	High = rpio.High
)

//Setup initialize memory buffers for GPIO
func Setup() error {
	return rpio.Open()
}

//Pin minimal interface for a GPIO pin
type Pin interface {
	Input()
	Output()
	High()
	Low()
	Read() rpio.State
}

//Open open a handler for a specific GPIO pin
func Open(pin int) Pin {
	return rpio.Pin(pin)
}

//WaitChange wait until the pin reliably reads `mode`, returning the elapsed duration
//If `timeout` ms elapse in the meantime, returns an error instead
func WaitChange(pin Pin, mode rpio.State, timeout time.Duration) (time.Duration, error) {
	start := time.Now()

	for {
		elapsed := time.Now().Sub(start)

		if elapsed > timeout {
			return 0 * time.Microsecond, errors.New("timeout")
		}

		a := pin.Read()
		b := pin.Read()
		c := pin.Read()

		if (a == b) && (b == c) && (c == mode) {
			return elapsed, nil
		}
	}
}
