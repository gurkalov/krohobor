package cli

import (
	"fmt"
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
	s3Storage := storage.AwsS3{
		Bucket: cfg.App.Catalog,
	}

	app.Flags = []cli.Flag {
		&cli.StringFlag{
			Name: "db",
			Value: "app",
			Usage: "database",
		},
		&cli.IntFlag{
			Name: "port",
			Value: 80,
			Usage: "Port",
		},
	}

	app.Commands = []*cli.Command{
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
			Name:    "complete",
			Aliases: []string{"c"},
			Usage:   "complete a task on the list",
			Action:  func(c *cli.Context) error {
				fmt.Println("completed task: ", c.Args().First())
				return nil
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
				{
					Name:  "dumpall",
					Usage: "backup all databases",
					Action: (actions.DbDumpAll{
						UseCase: usecases.NewDbDumpAll(dbPostgres, zipArchive, s3Storage),
					}).Action(cfg),
				},
			},
		},
	}


	return app
}
