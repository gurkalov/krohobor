package actions

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"krohobor/app/adapters/presenter/tsv"
	"krohobor/app/usecases"
)

type DumpList struct {
	UseCase usecases.DumpListInterface
}

func (d DumpList) Action(c *cli.Context) error {
	request := usecases.DumpListRequest{}

	resp, err := d.UseCase.Execute(request)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	pr := tsv.Presenter{}

	return pr.Print(resp.List)
}
