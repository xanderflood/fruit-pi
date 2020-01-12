package main

import (
	"fmt"
	"log"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/stianeikeland/go-rpio"
	"github.com/xanderflood/fruit-pi/lib/tools"
	"github.com/xanderflood/fruit-pi/pkg/gpio"
)

//Timeout defualt timeout for reading sequences
const Timeout = time.Second

var opts struct {
	Pin int `long:"pin" env:"PIN" default:"22"`
}

//GPIOStates GPIO state names
var GPIOStates = map[gpio.State]string{
	gpio.Low:  "low",
	gpio.High: "high",
}

func main() {
	fmt.Println("-- Starting Fruit-Pi Client --")

	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	if err := gpio.Setup(); err != nil {
		log.Fatal(err)
	}
	pin := rpio.Pin(opts.Pin)

	logger := tools.NewStdoutLogger(tools.LogLevelDebug, "send")

	for {
		logger.Infof("true")
		gpio.Set(pin, true)
		time.Sleep(time.Second)

		logger.Infof("false")
		gpio.Set(pin, false)
		time.Sleep(time.Second)
	}
}
