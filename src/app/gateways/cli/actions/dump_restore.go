package actions

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"krohobor/app/adapters/presenter/tsv"
	"krohobor/app/usecases"
)

type DumpRestore struct {
	UseCase usecases.DumpRestoreInterface
}

func (d DumpRestore) Action(c *cli.Context) error {
	request := usecases.DumpRestoreRequest{
		Name:     c.String("dbname"),
		Filename: c.String("name"),
	}

	resp, err := d.UseCase.Execute(request)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	pr := tsv.Presenter{}
	return pr.Print(resp)
}
