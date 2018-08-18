package am2301

import (
	"errors"
	"time"

	"github.com/xanderflood/fruit-pi/pkg/gpio"
)

////////////////
// this is a golang port of the C library found at:
// https://github.com/kporembinski/DHT21-AM2301/blob/master/am2301.c
////////////////

//State state of an AM2301 sensor
type State struct {
	RH   float32
	Temp float32
}

//AM2301 standard implementation of a Monitor
type AM2301 struct {
	pin      gpio.Pin
	interval time.Duration
	mode     int // should always be 1 - something to do with restarting
}

//Start start a monitor with the given handler hook
func (am *AM2301) Start(h func(State)) {
	go am.watch(h)
}

func (am *AM2301) watch(h func(State)) {
	for {
		time.Sleep(am.interval)

		err := am.Request()
		if err != nil {
			//TODO log these somewhere
			continue
		}

		vals, err := am.Read()
		if err != nil {
			//TODO log these somewhere
			continue
		}

		//parse the results
		state, err := Parse(vals)
		if err != nil {
			//invalid checksum
			//TODO log these somewhere
			continue
		}

		if !state.Valid() {
			//TODO log these somewhere
			continue
		}

		h(state)
	}
}

//Request send the signal to request a new measurement
func (am *AM2301) Request() error {
	// Leave it high for a while
	am.pin.Output()
	am.pin.High()
	time.Sleep(100 * time.Microsecond)

	// Set it low to give the start signal
	am.pin.Low()
	time.Sleep(1000 * time.Microsecond)

	// Now set the pin high to let the sensor start communicating
	am.pin.High()
	am.pin.Input()
	if _, err := gpio.WaitChange(am.pin, gpio.High, 100*time.Microsecond); err != nil {
		return errors.New("unexpected sequence")
	}

	// Wait for ACK
	if _, err := gpio.WaitChange(am.pin, gpio.Low, 100*time.Microsecond); err != nil {
		return errors.New("unexpected sequence")
	}
	if _, err := gpio.WaitChange(am.pin, gpio.High, 100*time.Microsecond); err != nil {
		return errors.New("unexpected sequence")
	}

	// When restarting, it looks like this look for start bit is not needed
	if am.mode != 0 {
		// Wait for the start bit
		if _, err := gpio.WaitChange(am.pin, gpio.Low, 200*time.Microsecond); err != nil {
			return errors.New("unexpected sequence")
		}

		if _, err := gpio.WaitChange(am.pin, gpio.High, 200*time.Microsecond); err != nil {
			return errors.New("unexpected sequence")
		}
	}

	return nil
}

//Read read a 5-byte sequence from the pin
func (am *AM2301) Read() ([5]byte, error) {
	var vals [5]byte
	for i := 0; i < 5; i++ {
		for j := 7; j >= 0; j-- {
			val, err := gpio.WaitChange(am.pin, gpio.Low, 500*time.Microsecond)
			if err != nil {
				return [5]byte{}, errors.New("unexpected signal")
			}

			if val >= 50*time.Microsecond {
				vals[i] = vals[i] | (1 << uint(j))
			}

			_, err = gpio.WaitChange(am.pin, gpio.High, 500*time.Microsecond)
			if err != nil {
				return [5]byte{}, errors.New("unexpected signal")
			}
		}
	}

	am.pin.Output()
	am.pin.High()

	return vals, nil
}

//Parse parse a State from a 5-byte input stream
func Parse(vals [5]byte) (State, error) {
	//TODO test this using the examples on page 5 of
	//https://kropochev.com/downloads/humidity/AM2301.pdf

	// Verify checksum
	if vals[0]+vals[1]+vals[2]+vals[3] != vals[4] {
		//TODO log these somewhere
		return State{}, errors.New("invalid checksum")
	}

	tempSign := 1
	if (vals[2] >> 7) != 0 {
		//turn off the sign bit and set the sign
		vals[2] ^= (1 << 7)
		tempSign = -1
	}

	return State{
		RH:   float32((int(vals[0])<<8)|int(vals[1])) / 10.0,
		Temp: float32(tempSign*((int(vals[2])<<8)|int(vals[3]))) / 10.0,
	}, nil
}

//Valid check that values are within the specifed range
func (s State) Valid() bool {
	return (s.RH <= 100.0) && (s.RH >= 0.0) && (s.Temp <= 80.0) && (s.Temp <= 40.0)
}
