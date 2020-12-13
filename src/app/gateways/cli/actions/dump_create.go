package actions

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"krohobor/app/adapters/config"
	"krohobor/app/usecases"
	"strings"
)

type DumpCreate struct {
	UseCase usecases.DumpCreateInterface
}

func (d DumpCreate) Action(cfg config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		dbs := c.String("db")

		var dbNames []string
		if dbs != "" {
			dbNames = strings.Split(dbs, ",")
		} else {
			dbs = "all"
		}

		filename := fmt.Sprintf("/tmp/backup/%s.sql", dbs)

		request := usecases.DumpCreateRequest{
			DbNames: dbNames,
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
