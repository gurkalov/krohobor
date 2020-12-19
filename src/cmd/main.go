package main

import (
	"flag"
	"krohobor/app/adapters/config"
	"krohobor/app/gateways/cli"
	"log"
	"os"
)

func main() {
	configFile := flag.String("config", ".env", "a string")
	cfg := config.Load(*configFile)

	app := cli.App(cfg)
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
