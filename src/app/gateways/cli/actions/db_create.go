package actions

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"krohobor/app/usecases"
)

type DbCreate struct {
	UseCase usecases.DbCreateInterface
}

func (d DbCreate) Action(c *cli.Context) error {
	request := usecases.DbCreateRequest{
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
