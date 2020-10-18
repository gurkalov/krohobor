package main

import (
	"krohobor/app/adapters/config"
	"krohobor/app/gateways/cli"
	"log"
	"os"
)

func main() {
	cfg := config.Load()

	app := cli.App(cfg)
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
