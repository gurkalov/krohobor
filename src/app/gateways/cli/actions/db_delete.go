package actions

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"krohobor/app/adapters/presenter/tsv"
	"krohobor/app/usecases"
)

type DbDelete struct {
	UseCase usecases.DbDeleteInterface
}

func (d DbDelete) Action(c *cli.Context) error {
	request := usecases.DbDeleteRequest{
		Name:  c.String("dbname"),
		Force: c.Bool("force"),
	}
	resp, err := d.UseCase.Execute(request)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	pr := tsv.Presenter{}
	return pr.Print(resp)
}
