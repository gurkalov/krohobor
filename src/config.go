package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port int
	Catalog string
	Password string
}

var cfg Config

func InitConfig() {
	configFile := flag.String("config", ".env", "a string")
	if _, err := toml.DecodeFile(*configFile, &cfg); err != nil {
		port, _ := strconv.Atoi(os.Getenv("APP_PORT"))
		cfg.Port = port
		cfg.Password = os.Getenv("APP_PASSWORD")
		cfg.Catalog = os.Getenv("APP_CATALOG")

		pgpass := []byte(os.Getenv("PGHOST") +
			":" + os.Getenv("PGPORT") +
			":postgres" +
			":" + os.Getenv("PGUSER") +
			":" + os.Getenv("PGPASSWORD"))
		if err := ioutil.WriteFile("/root/.pgpass", pgpass, 0600); err != nil {
			log.Fatalln(err)
		}
	}
}
