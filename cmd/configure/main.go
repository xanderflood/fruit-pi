package main

import (
	"fmt"
	"log"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/xanderflood/fruit-pi/pkg/config"
)

var opts struct {
	Config string `long:"config-path" env:"CONFIG_PATH" default:"./config/env.json"`
}

func main() {
	fmt.Println("-- Fruit-Pi Configuration Update --")

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

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.Load(f)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Applying updates...")
	err = cfg.Apply(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Saving changes...")
	f, err = os.OpenFile(opts.Config, os.O_RDWR, 0755)
	err = cfg.Save(f)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Configuration Successful")
}
