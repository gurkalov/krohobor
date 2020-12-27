package database

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"krohobor/app/adapters/config"
)

func Config(name string, cfg config.Config) config.DatabaseConfig {
	var dbConfig config.DatabaseConfig
	for _, v := range cfg.Databases {
		if v.Name == name {
			dbConfig = v
			break
		}
	}

	return dbConfig
}

func Impl(dbConfig config.DatabaseConfig) (Interface, error) {
	switch dbConfig.Driver {
	case "postgres": {
		var conf config.PostgresConfig
		if err := mapstructure.Decode(dbConfig.Options, &conf); err != nil {
			return nil, errors.New(fmt.Sprintf("Database %s not found", dbConfig.Name))
		}
		return NewPostgres(conf), nil
	}
	case "memory": {
		return NewMemory(), nil
	}
	}

	return nil, errors.New(fmt.Sprintf("Database %s driver %s error", dbConfig.Name, dbConfig.Driver))
}
