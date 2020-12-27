package actions

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"krohobor/app/usecases"
)

type DumpRestore struct {
	UseCase usecases.DumpRestoreInterface
}

func (d DumpRestore) Action(c *cli.Context) error {
	request := usecases.DumpRestoreRequest{
		Name: c.String("name"),
		Filename: c.String("name"),
	}

	resp, err := d.UseCase.Execute(request)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(resp)

	return nil
}
