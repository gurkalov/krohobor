package storage

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"krohobor/app/adapters/archive"
	"krohobor/app/adapters/config"
)

func Config(name string, cfg config.Config) (config.StorageConfig, error) {
	var storageConfig config.StorageConfig
	for _, v := range cfg.Storages {
		if v.Name == name {
			storageConfig = v
			return storageConfig, nil
		}
	}

	return storageConfig, errors.New(fmt.Sprintf("Storage %s not found", name))
}

func Impl(storageConfig config.StorageConfig, arch archive.Interface) (Interface, error) {
	switch storageConfig.Driver {
	case "s3": {
		var conf config.AwsS3Config
		if err := mapstructure.Decode(storageConfig.Options, &conf); err != nil {
			return nil, errors.New(fmt.Sprintf("Storage %s error: %v", storageConfig.Name, err))
		}
		return NewAwsS3(conf, arch), nil
	}
	case "file": {
		var conf config.FileConfig
		if err := mapstructure.Decode(storageConfig.Options, &conf); err != nil {
			return nil, errors.New(fmt.Sprintf("Storage %s error: %v", storageConfig.Name, err))
		}
		return NewFile(conf.Catalog, arch), nil
	}
	}

	return nil, errors.New(fmt.Sprintf("Storage %s driver %s not found", storageConfig.Name, storageConfig.Driver))
}
