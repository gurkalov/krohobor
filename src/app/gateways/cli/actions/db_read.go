package actions

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"krohobor/app/usecases"
)

type DbRead struct {
	UseCase usecases.DbReadInterface
}

func (d DbRead) Action(c *cli.Context) error {
	request := usecases.DbReadRequest{
		Name: c.String("dbname"),
	}
	resp, err := d.UseCase.Execute(request)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(resp)

	return nil
}
