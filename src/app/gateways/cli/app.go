package cli

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/urfave/cli/v2"
	"krohobor/app/adapters/archive"
	"krohobor/app/adapters/config"
	"krohobor/app/adapters/database"
	"krohobor/app/adapters/storage"
	"krohobor/app/gateways/cli/actions"
	"krohobor/app/gateways/httpserver"
	"krohobor/app/usecases"
	"net/http"
	"strconv"
)

func InitDatabase(cfg config.Config) (database.Interface, error) {
	var dbConfig config.DatabaseConfig
	for _, v := range cfg.Databases {
		if v.Name == cfg.App.Database {
			dbConfig = v
			break
		}
	}

	if dbConfig.Name == "" {
		return nil, errors.New(fmt.Sprintf("Database %s not found", cfg.App.Database))
	}

	switch dbConfig.Driver {
	case "postgres": {
		var conf config.PostgresConfig
		if err := mapstructure.Decode(dbConfig.Options, &conf); err != nil {
			return nil, errors.New(fmt.Sprintf("Database %s not found", cfg.App.Database))
		}
		return database.NewPostgres(conf), nil
	}
	case "memory": {
		return database.NewMemory(), nil
	}
	}

	return nil, errors.New(fmt.Sprintf("Database %s driver %s error", cfg.App.Database, dbConfig.Driver))
}

func InitStorage(cfg config.Config, arch archive.Interface) (storage.Interface, error) {
	var storageConfig config.StorageConfig
	for _, v := range cfg.Storages {
		if v.Name == cfg.App.Storage {
			storageConfig = v
			break
		}
	}

	if storageConfig.Name == "" {
		return nil, errors.New(fmt.Sprintf("Storage %s not found", cfg.App.Storage))
	}

	switch storageConfig.Driver {
	case "s3": {
		var conf config.AwsS3Config
		if err := mapstructure.Decode(storageConfig.Options, &conf); err != nil {
			return nil, errors.New(fmt.Sprintf("Storage %s error: %v", cfg.App.Storage, err))
		}
		return storage.NewAwsS3(conf, arch), nil
	}
	case "file": {
		var conf config.FileConfig
		if err := mapstructure.Decode(storageConfig.Options, &conf); err != nil {
			return nil, errors.New(fmt.Sprintf("Storage %s error: %v", cfg.App.Storage, err))
		}
		return storage.NewFile(conf.Catalog, arch), nil
	}
	}

	return nil, errors.New(fmt.Sprintf("Storage %s driver %s error", cfg.App.Storage, storageConfig.Driver))
}

func App(cfg config.Config) *cli.App {
	app := &cli.App{}

	db, err := InitDatabase(cfg)
	if err != nil {
		panic(err)
	}

	zipArchive := archive.Zip{
		Password: cfg.App.Password,
	}

	store, err := InitStorage(cfg, zipArchive)
	if err != nil {
		panic(err)
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name: "db",
			Value: "",
			Usage: "database",
		},
		&cli.StringFlag{
			Name: "name",
			Value: "",
			Usage: "Backup name",
		},
		&cli.StringFlag{
			Name: "target",
			Value: "",
			Usage: "Target host",
		},
		&cli.IntFlag{
			Name: "port",
			Value: 80,
			Usage: "Port",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:    "status",
			Aliases: []string{},
			Usage:   "Status info",
			Action: (actions.Status{
				UseCase: usecases.NewStatus(cfg, db, store),
			}).Action(cfg),
		},
		{
			Name:    "httpserver",
			Aliases: []string{"serve"},
			Usage:   "HTTP Server",
			Action:  func(c *cli.Context) error {
				router := httpserver.Router(cfg)
				return http.ListenAndServe(":" + strconv.Itoa(c.Int("port")), router)
			},
		},
		{
			Name:        "db",
			Aliases:     []string{"db"},
			Usage:       "options for task templates",
			Subcommands: []*cli.Command{
				{
					Name:  "list",
					Usage: "list of databases",
					Action: (actions.DbList{
						UseCase: usecases.NewDbList(db),
					}).Action(cfg),
				},
				{
					Name:  "read",
					Usage: "read a database",
					Action: (actions.DbRead{
						UseCase: usecases.NewDbRead(db),
					}).Action(cfg),
				},
			},
		},
		{
			Name:        "dump",
			Aliases:     []string{},
			Usage:       "Work with dumps",
			Subcommands: []*cli.Command{
				{
					Name:  "create",
					Usage: "create dump",
					Action: (actions.DumpCreate{
						UseCase: usecases.NewDumpCreate(db, store),
					}).Action(cfg),
				},
				{
					Name:  "list",
					Usage: "list of dumps",
					Action: (actions.DumpList{
						UseCase: usecases.NewDumpList(store),
					}).Action(cfg),
				},
				{
					Name:  "restore",
					Usage: "restore dump",
					Action: (actions.DumpRestore{
						UseCase: usecases.NewDumpRestore(db, store),
					}).Action(cfg),
				},
				{
					Name:  "delete",
					Usage: "delete backup",
					Action: (actions.DumpDelete{
						UseCase: usecases.NewDumpDelete(store),
					}).Action(cfg),
				},
			},
		},
	}

	return app
}
