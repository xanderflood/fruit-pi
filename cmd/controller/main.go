package main

import (
	"context"
	"log"
	"net/http"
	"time"

	flags "github.com/jessevdk/go-flags"
	rpio "github.com/stianeikeland/go-rpio"

	"github.com/xanderflood/fruit-pi-server/lib/api"
	"github.com/xanderflood/fruit-pi/lib/tools"
	"github.com/xanderflood/fruit-pi/pkg/device"
)

var opts struct {
	Host                     string `long:"fruit-pi-host" env:"FRUIT_PI_HOST" required:"true"`
	Token                    string `long:"fruit-pi-token" env:"FRUIT_PI_TOKEN" required:"true"`
	StatePersistenceLocation string `long:"state-persistence-location" env:"STATE_PERSISTENCE_LOCATION" default:"./config.json"`
	SkipGPIOInitialization   bool   `long:"skip-gpio-initialization" env:"SKIP_GPIO_INITIALIZATION" optional:"true"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	if !opts.SkipGPIOInitialization {
		err = rpio.Open()
		if err != nil {
			log.Fatal(err)
		}
	}

	logger := tools.NewStdoutLogger(tools.LogLevelDebug, "fruit-pi")
	client := api.NewDefaultClient(
		opts.Host,
		http.DefaultTransport,
		opts.Token,
	)
	dvc := device.New(client, logger)

	err = dvc.Start(
		context.Background(),
		opts.StatePersistenceLocation,
		5*time.Second,
		10*time.Second,
	)
	if err != nil {
		log.Fatal(err)
	}

	//permenantly pause this goroutine
	<-make(chan interface{})
	return
}
