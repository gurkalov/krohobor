package main

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port int
	Password string
}

var cfg Config

func InitConfig() {
	if _, err := toml.DecodeFile(".env", &cfg); err != nil {
		port, _ := strconv.Atoi(os.Getenv("APP_PORT"))
		cfg.Port = port

		cfg.Password = os.Getenv("APP_PASSWORD")

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
