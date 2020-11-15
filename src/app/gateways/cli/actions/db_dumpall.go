package actions

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"krohobor/app/adapters/config"
	"krohobor/app/usecases"
)

type DbDumpAll struct {
	UseCase usecases.DbDumpAllInterface
}

func (d DbDumpAll) Action(cfg config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		filename := "/tmp/backup/all.sql"

		request := usecases.DbDumpAllRequest{
			Filename: filename,
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
