package actions

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"krohobor/app/usecases"
)

type DbList struct {
	UseCase usecases.DbListInterface
}

func (d DbList) Action(c *cli.Context) error {
	request := usecases.DbListRequest{}
	resp, err := d.UseCase.Execute(request)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(resp)

	return nil
}
