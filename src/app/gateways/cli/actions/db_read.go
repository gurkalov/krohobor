package actions

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"krohobor/app/adapters/config"
	"krohobor/app/usecases"
)

type DbRead struct {
	UseCase usecases.DbReadInterface
}

func (d DbRead) Action(cfg config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		request := usecases.DbReadRequest{
			Name: c.String("db"),
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
