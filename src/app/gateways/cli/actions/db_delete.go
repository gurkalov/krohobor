package actions

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"krohobor/app/usecases"
)

type DbDelete struct {
	UseCase usecases.DbDeleteInterface
}

func (d DbDelete) Action(c *cli.Context) error {
	request := usecases.DbDeleteRequest{
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
