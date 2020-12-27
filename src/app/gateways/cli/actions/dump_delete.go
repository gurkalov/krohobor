package actions

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"krohobor/app/usecases"
)

type DumpDelete struct {
	UseCase usecases.DumpDeleteInterface
}

func (d DumpDelete) Action(c *cli.Context) error {
	request := usecases.DumpDeleteRequest{
		Name: c.String("name"),
	}

	resp, err := d.UseCase.Execute(request)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(resp)

	return nil
}
