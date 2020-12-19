package config

import (
	"github.com/BurntSushi/toml"
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

func Load(filename string) Config {
	cfg := Config{}

	if _, err := toml.DecodeFile(filename, &cfg); err != nil {
		port, _ := strconv.Atoi(os.Getenv("APP_PORT"))
		cfg.App.Port = port
		cfg.App.Password = os.Getenv("APP_PASSWORD")
		cfg.App.Catalog = os.Getenv("APP_CATALOG")

		cfg.Postgres.Host = os.Getenv("PGHOST")
		cfg.Postgres.Port = os.Getenv("PGPORT")
		cfg.Postgres.DB = os.Getenv("PGDB")
		cfg.Postgres.User = os.Getenv("PGUSER")
		cfg.Postgres.Password = os.Getenv("PGPASSWORD")
	}

	return cfg
}

func LoadMock() Config {
	cfg := Config{}

	return cfg
}
