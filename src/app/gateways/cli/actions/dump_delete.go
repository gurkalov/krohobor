package actions

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"krohobor/app/adapters/presenter/tsv"
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

	pr := tsv.Presenter{}
	return pr.Print(resp)
}
