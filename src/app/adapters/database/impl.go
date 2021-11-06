package database

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"krohobor/app/adapters/config"
)

func Config(name string, cfg config.Config) (config.DatabaseConfig, error) {
	var dbConfig config.DatabaseConfig
	for _, v := range cfg.Databases {
		if v.Name == name {
			dbConfig = v
			return dbConfig, nil
		}
	}

	return dbConfig, fmt.Errorf("Database %s not found", name)
}

func Impl(dbConfig config.DatabaseConfig) (Interface, error) {
	switch dbConfig.Driver {
	case "postgres":
		{
			var conf config.PostgresConfig
			if err := mapstructure.Decode(dbConfig.Options, &conf); err != nil {
				return nil, fmt.Errorf("Database %s not found", dbConfig.Name)
			}
			return NewPostgres(conf), nil
		}
	case "memory":
		{
			return NewMemory(), nil
		}
	}

	return nil, fmt.Errorf("Database %s driver %s not found", dbConfig.Name, dbConfig.Driver)
}
