package main

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port          int
}

var cfg Config

func InitConfig() {
	if _, err := toml.DecodeFile(".env", &cfg); err != nil {
		port, _ := strconv.Atoi(os.Getenv("APP_PORT"))
		cfg.Port = port

		d1 := []byte("postgres:5432:postgres:postgres:secret")
		if err := ioutil.WriteFile("/root/.pgpass", d1, 0600); err != nil {
			log.Fatalln(err)
		}
	}
}
