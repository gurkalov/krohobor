package actions

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"krohobor/app/usecases"
)

type Status struct {
	UseCase usecases.StatusInterface
}

func (d Status) Action(c *cli.Context) error {
	request := usecases.StatusRequest{}
	resp, err := d.UseCase.Execute(request)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(resp)

	return nil
}
