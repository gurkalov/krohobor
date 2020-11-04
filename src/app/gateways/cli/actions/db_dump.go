package actions

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"krohobor/app/adapters/config"
	"krohobor/app/usecases"
)

type DbDump struct {
	UseCase usecases.DbDumpInterface
}

func (d DbDump) Action(cfg config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		name := c.String("db")

		filename := fmt.Sprintf("/tmp/backup/%s.sql", name)
		dirname := "/tmp/backup"
		archname := fmt.Sprintf("/tmp/backup_%s.zip", name)

		request := usecases.DbDumpRequest{
			Name: name,
			Filename: filename,
			Dirname: dirname,
			Archname: archname,
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
