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
	"time"
)

func UseDatabase(name string, cfg config.Config) (database.Interface, error) {
	if name == "" {
		name = cfg.App.Database
	}

	dbCfg, err := database.Config(name, cfg)
	if err != nil {
		return nil, err
	}

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
		Dir:      cfg.App.Dir,
		Password: cfg.App.Password,
	}

	storeCfg, err := storage.Config(name, cfg)
	if err != nil {
		return nil, err
	}

	store, err := storage.Impl(storeCfg, zipArchive)
	if err != nil {
		return store, err
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
		&cli.BoolFlag{
			Name:  "force",
			Value: false,
			Usage: "Force flag",
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
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:  "human",
							Usage: "show human-readable size",
						},
					},
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
						// Track the progress of loading the dump to S3 storage.
						if s3store, ok := store.(*storage.AwsS3); ok {
							s3store.SetProgress(actions.ProgressPrinter("Uploading"))
						}
						return (actions.DumpCreate{
							UseCase: usecases.NewDumpCreate(db, store),
						}).Action(c, cfg.App)
					},
				},
				{
					Name:  "list",
					Usage: "list of dumps",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:  "human",
							Usage: "show human-readable size",
						},
					},
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
						start := time.Now()
						actions.LogRestoreStart()

						// Track the progress of downloading the dump from S3 storage
						// and log the download and unarchive between the steps.
						if s3store, ok := store.(*storage.AwsS3); ok {
							s3store.SetProgress(actions.ProgressPrinter("Downloading"))
							s3store.SetOnDownloaded(actions.LogDownloadComplete)
							s3store.SetOnUnarchived(actions.LogUnarchiveComplete)
						}
						// Track the progress of restoring the dump into the database.
						if pgdb, ok := db.(*database.Postgres); ok {
							pgdb.SetProgress(actions.ProgressPrinter("Restoring"))
						}

						err := (actions.DumpRestore{
							UseCase: usecases.NewDumpRestore(db, store),
						}).Action(c)

						actions.LogRestoreFinished(time.Since(start))
						return err
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
