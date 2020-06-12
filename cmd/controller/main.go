package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	flags "github.com/jessevdk/go-flags"
	rpio "github.com/stianeikeland/go-rpio"

	"github.com/xanderflood/fruit-pi-server/lib/api"
	"github.com/xanderflood/fruit-pi/pkg/device"
)

var opts struct {
	Host                     *string `long:"fruit-pi-host" env:"FRUIT_PI_HOST"`
	Token                    *string `long:"fruit-pi-token" env:"FRUIT_PI_TOKEN"`
	StatePersistenceLocation string  `long:"state-persistence-location" env:"STATE_PERSISTENCE_LOCATION" default:"./config.json"`

	SkipGPIOInitialization bool    `long:"skip-gpio-initialization" env:"SKIP_GPIO_INITIALIZATION" optional:"true"`
	Base64DeviceConfig     *string `long:"base64-device-config" env:"BASE64_DEVICE_CONFIG"`
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

	var dvcResp api.Device
	if opts.Base64DeviceConfig == nil {
		if opts.Host == nil || opts.Token == nil {
			log.Fatal("--fruit-pi-host and --fruit-pi-token are required arguments")
		}

		client := api.NewDefaultClient(
			*opts.Host,
			http.DefaultTransport,
			*opts.Token,
		)

		dvcResp, err = client.GetDeviceConfig(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	} else {
		b64Resp, err := base64.StdEncoding.DecodeString(*opts.Base64DeviceConfig)
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(b64Resp, &dvcResp)
		if err != nil {
			log.Fatal(err)
		}
	}

	spew.Dump(*dvcResp.Config)

	var cfg device.Config
	err = json.Unmarshal(*dvcResp.Config, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	graph, err := cfg.BuildGraph()
	if err != nil {
		log.Fatal(err)
	}

	graph.Start(context.Background())

	// permanently pause the main goroutine
	<-make(chan interface{})
	return
}
