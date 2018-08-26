package main

import (
	"fmt"
	"log"
	"os"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/xanderflood/fruit-pi/pkg/chamber"
	"github.com/xanderflood/fruit-pi/pkg/config"
)

var opts struct {
	Config string `long:"config-path" env:"CONFIG_PATH" default:"./config/env.json"`
}

func main() {
	fmt.Println("Loading configuration...")
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(opts.Config, os.O_RDWR, 0755)
	if err != nil { //TODO: if file_not_Found, move on
		if os.IsNotExist(err) {
			fmt.Println(" - config does not exist, so creating it")
		} else {
			log.Fatal(err)
		}
	}

	cfg, err := config.Load(f)
	if err != nil {
		log.Fatal(err)
	}

	c := chamber.New(
		cfg.FanPin,
		cfg.HumPin,
		cfg.SensorPin,
		cfg.Strategy.Object)

	fmt.Println("Starting chamber")
	for {
		time.Sleep(cfg.Interval.Duration)

		fmt.Println("Refreshing chamber")
		cState, sState, err := c.Refresh()
		if err != nil {
			fmt.Println("Refresh failed: %s", err.Error())
		}

		fmt.Println("RlH:\t%v%", sState.RH)
		fmt.Println("Tmp:\t%vC", sState.Temp)
		fmt.Println("Hum:\t%t", cState.Hum)
		fmt.Println("Fan:\t%t", cState.Fan)
	}
}
