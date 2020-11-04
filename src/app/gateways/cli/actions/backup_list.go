package actions

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"krohobor/app/adapters/config"
	"krohobor/app/usecases"
)

type BackupList struct {
	UseCase usecases.BackupListInterface
}

func (d BackupList) Action(cfg config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		request := usecases.BackupListRequest{}

		resp, err := d.UseCase.Execute(request)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		fmt.Println(resp)

		return nil
	}
}
