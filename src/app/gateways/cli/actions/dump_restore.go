package actions

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"krohobor/app/adapters/config"
	"krohobor/app/usecases"
)

type DumpRestore struct {
	UseCase usecases.DumpRestoreInterface
}

func (d DumpRestore) Action(cfg config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		request := usecases.DumpRestoreRequest{
			Name: c.String("name"),
			Filename: c.String("name"),
			Target: c.String("target"),
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
