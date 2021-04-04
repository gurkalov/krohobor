package cli

import (
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

func UseDatabase(name string, cfg config.Config) (database.Interface, error) {
	if name == "" {
		name = cfg.App.Database
	}

	dbCfg := database.Config(name, cfg)
	db, err := database.Impl(dbCfg)
	if err != nil {
		return db, err
	}

	return db, nil
}

func UseStorage(name string, cfg config.Config) (storage.Interface, error) {
	if name == "" {
		name = cfg.App.Storage
	}

	zipArchive := archive.Zip{
		Password: cfg.App.Password,
	}

	storeCfg := storage.Config(name, cfg)
	store, err := storage.Impl(storeCfg, zipArchive)
	if err != nil {
		return store, nil
	}

	return store, nil
}

func App(cfg config.Config) *cli.App {
	app := &cli.App{}

	var db database.Interface
	var store storage.Interface
	initDeps := func(c *cli.Context) error {
		var err error
		db, err = UseDatabase(c.String("database"), cfg)
		if err != nil {
			return err
		}

		store, err = UseStorage(c.String("storage"), cfg)
		if err != nil {
			return err
		}

		return nil
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "dbname",
			Value: "",
			Usage: "Database name",
		},
		&cli.StringFlag{
			Name:  "name",
			Value: "",
			Usage: "Backup name",
		},
		&cli.StringFlag{
			Name:  "database",
			Value: "",
			Usage: "Use database",
		},
		&cli.StringFlag{
			Name:  "storage",
			Value: "",
			Usage: "Use storage",
		},
		&cli.IntFlag{
			Name:  "port",
			Value: 80,
			Usage: "Port",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:    "status",
			Aliases: []string{},
			Usage:   "Status info",
			Before:  initDeps,
			Action: func(c *cli.Context) error {
				return (actions.Status{
					UseCase: usecases.NewStatus(cfg, db, store),
				}).Action(c)
			},
		},
		{
			Name:    "httpserver",
			Aliases: []string{"serve"},
			Usage:   "HTTP Server",
			Action: func(c *cli.Context) error {
				router := httpserver.Router(cfg)
				return http.ListenAndServe(":"+strconv.Itoa(c.Int("port")), router)
			},
		},
		{
			Name:    "db",
			Aliases: []string{"db"},
			Usage:   "options for task templates",
			Before:  initDeps,
			Subcommands: []*cli.Command{
				{
					Name:  "list",
					Usage: "list of databases",
					Action: func(c *cli.Context) error {
						return (actions.DbList{
							UseCase: usecases.NewDbList(db),
						}).Action(c)
					},
				},
				{
					Name:  "read",
					Usage: "read a database",
					Action: func(c *cli.Context) error {
						return (actions.DbRead{
							UseCase: usecases.NewDbRead(db),
						}).Action(c)
					},
				},
				{
					Name:  "create",
					Usage: "create a database",
					Action: func(c *cli.Context) error {
						return (actions.DbCreate{
							UseCase: usecases.NewDbCreate(db),
						}).Action(c)
					},
				},
				{
					Name:  "delete",
					Usage: "delete a database",
					Action: func(c *cli.Context) error {
						return (actions.DbDelete{
							UseCase: usecases.NewDbDelete(db),
						}).Action(c)
					},
				},
			},
		},
		{
			Name:    "dump",
			Aliases: []string{},
			Usage:   "Work with dumps",
			Before:  initDeps,
			Subcommands: []*cli.Command{
				{
					Name:  "create",
					Usage: "create dump",
					Action: func(c *cli.Context) error {
						return (actions.DumpCreate{
							UseCase: usecases.NewDumpCreate(db, store),
						}).Action(c)
					},
				},
				{
					Name:  "list",
					Usage: "list of dumps",
					Action: func(c *cli.Context) error {
						return (actions.DumpList{
							UseCase: usecases.NewDumpList(store),
						}).Action(c)
					},
				},
				{
					Name:  "restore",
					Usage: "restore dump",
					Action: func(c *cli.Context) error {
						return (actions.DumpRestore{
							UseCase: usecases.NewDumpRestore(db, store),
						}).Action(c)
					},
				},
				{
					Name:  "delete",
					Usage: "delete backup",
					Action: func(c *cli.Context) error {
						return (actions.DumpDelete{
							UseCase: usecases.NewDumpDelete(store),
						}).Action(c)
					},
				},
			},
		},
	}

	return app
}
