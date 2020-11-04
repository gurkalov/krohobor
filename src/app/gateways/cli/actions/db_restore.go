package actions

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"krohobor/app/adapters/config"
	"krohobor/app/usecases"
)

type DbRestore struct {
	UseCase usecases.DbRestoreInterface
}

func (d DbRestore) Action(cfg config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		filename := "/tmp/download_backup/all.sql"
		dirname := "/tmp/download_backup"
		archname := "/tmp/download_backup.zip"

		request := usecases.DbRestoreRequest{
			Filename: filename,
			Name: c.String("name"),
			Archfile : archname,
			Archdir: dirname,
		}

		resp, err := d.UseCase.Execute(request)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		fmt.Println(resp)

		return nil
	}
}
