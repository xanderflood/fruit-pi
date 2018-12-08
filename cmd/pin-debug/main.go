package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	flags "github.com/jessevdk/go-flags"
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
	pin := gpio.New(opts.Pin)

	logger := tools.NewStdoutLogger(tools.LogLevelDebug, "send")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("> ")
		scanner.Scan()

		seq, err := gpio.ToSequence(scanner.Text())
		if err != nil {
			logger.Error(err)
			continue
		}
		logger.Infof("EXECUTING:", seq.String())

		gpio.Execute(pin, seq)
		response, err := gpio.Monitor(pin, seq.NextState(), Timeout)
		if err != nil {
			logger.Error(err)
			continue
		}

		logger.Infof("GOT RESPONSE:")
		logger.Infof(response.String())
	}
}
