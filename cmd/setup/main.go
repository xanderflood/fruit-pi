package main

import (
	"fmt"
	"log"
	"os"

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
		nil /*TODO: strategy*/)

	err = c.Setup()
	if err != nil {
		log.Fatal(err)
	}
}
