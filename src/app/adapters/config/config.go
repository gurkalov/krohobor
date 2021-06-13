package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type Config struct {
	App       AppConfig
	Databases []DatabaseConfig
	Storages  []StorageConfig
}

type AppConfig struct {
	Port     int
	Dir      string
	Format   string
	Password string
	Database string
	Storage  string
}

type PostgresConfig struct {
	Host     string
	Port     int
	DB       string
	User     string
	Password string
}

type AwsS3Config struct {
	Catalog   string
	KeyId     string
	AccessKey string
	Region    string
}

type FileConfig struct {
	Catalog string
}

type DatabaseConfig struct {
	Name    string
	Driver  string
	Options map[string]interface{}
}

type StorageConfig struct {
	Name    string
	Driver  string
	Options map[string]interface{}
}

func Load(filename string) Config {
	cfg := Config{}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		port, _ := strconv.Atoi(os.Getenv("APP_PORT"))
		cfg.App.Port = port
	}

	return cfg
}

func LoadMock() Config {
	cfg := Config{}

	return cfg
}
