package actions

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"krohobor/app/usecases"
	"strings"
)

type DumpCreate struct {
	UseCase usecases.DumpCreateInterface
}

func (d DumpCreate) Action(c *cli.Context) error {
	dbs := c.String("dbname")

	var dbNames []string
	if dbs != "" {
		dbNames = strings.Split(dbs, ",")
	} else {
		dbs = "all"
	}

	filename := fmt.Sprintf("%s.sql", dbs)

	request := usecases.DumpCreateRequest{
		DbNames:  dbNames,
		Filename: filename,
	}

	resp, err := d.UseCase.Execute(request)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(resp)

	return nil
}
