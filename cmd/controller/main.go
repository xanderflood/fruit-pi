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
	Host  string
	Token string
}

func main() {
	fmt.Println("Loading configuration...")
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	err = rpio.Open()
	if err != nil {
		log.Fatal(err)
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
