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

func App(cfg config.Config) *cli.App {
	app := &cli.App{}

	dbPostgres := database.NewPostgres(cfg.Postgres)
	zipArchive := archive.Zip{
		Password: cfg.App.Password,
	}
	s3Storage := storage.NewAwsS3(cfg.App.Catalog, zipArchive)

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
				UseCase: usecases.NewStatus(cfg, dbPostgres, s3Storage),
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
						UseCase: usecases.NewDbList(dbPostgres),
					}).Action(cfg),
				},
				{
					Name:  "read",
					Usage: "read a database",
					Action: (actions.DbRead{
						UseCase: usecases.NewDbRead(dbPostgres),
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
						UseCase: usecases.NewDumpCreate(dbPostgres, s3Storage),
					}).Action(cfg),
				},
				{
					Name:  "list",
					Usage: "list of dumps",
					Action: (actions.DumpList{
						UseCase: usecases.NewDumpList(s3Storage),
					}).Action(cfg),
				},
				{
					Name:  "restore",
					Usage: "restore dump",
					Action: (actions.DumpRestore{
						UseCase: usecases.NewDumpRestore(dbPostgres, s3Storage),
					}).Action(cfg),
				},
				{
					Name:  "delete",
					Usage: "delete backup",
					Action: (actions.DumpDelete{
						UseCase: usecases.NewDumpDelete(s3Storage),
					}).Action(cfg),
				},
			},
		},
	}

	return app
}
