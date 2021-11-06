package main

import (
	"krohobor/app/adapters/config"
	"krohobor/app/gateways/cli"
	"log"
	"os"
)

func main() {
	configFile := os.Getenv("KROHOBOR_CONFIG")
	if configFile == "" {
		configFile = "config.yaml"
	}

	cfg := config.Load(configFile)

	app := cli.App(cfg)
	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
