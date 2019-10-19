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
	App AppConfig
	Postgres PostgresConfig
}

type AppConfig struct {
	Port     int
	Catalog  string
	Password string
}

type PostgresConfig struct {
	Host     string
	Port     string
	DB       string
	User     string
	Password string
}

var cfg Config

func InitConfig() {
	configFile := flag.String("config", ".env", "a string")
	if _, err := toml.DecodeFile(*configFile, &cfg); err != nil {
		port, _ := strconv.Atoi(os.Getenv("APP_PORT"))
		cfg.App.Port = port
		cfg.App.Password = os.Getenv("APP_PASSWORD")
		cfg.App.Catalog = os.Getenv("APP_CATALOG")

		cfg.Postgres.Host = os.Getenv("PGHOST")
		cfg.Postgres.Port = os.Getenv("PGPORT")
		cfg.Postgres.DB = "*"
		cfg.Postgres.User = os.Getenv("PGUSER")
		cfg.Postgres.Password = os.Getenv("PGPASSWORD")
	}

	pg := cfg.Postgres
	pgpass := []byte(pg.Host +
		":" + pg.Port +
		":*" +
		":" + pg.User +
		":" + pg.Password)
	if err := ioutil.WriteFile(".pgpass", pgpass, 0600); err != nil {
		log.Fatalln(err)
	}
}
