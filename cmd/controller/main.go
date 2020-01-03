package main

import (
	"context"
	"fmt"
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
	Host                   string `long:"fruit-pi-host" env:"FRUIT_PI_HOST" required:"true"`
	Token                  string `long:"fruit-pi-token" env:"FRUIT_PI_TOKEN" required:"true"`
	SkipGPIOInitialization bool   `long:"skip-gpio-initialization" env:"SKIP_GPIO_INITIALIZATION" optional:"true"`
}

func main() {
	fmt.Println("Loading configuration...")
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

	for {
		fmt.Println("Refreshing chamber")
		dvc.Refresh(context.Background())

		time.Sleep(time.Second)
	}
}
