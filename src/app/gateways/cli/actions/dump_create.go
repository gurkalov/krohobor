package actions

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"krohobor/app/adapters/config"
	"krohobor/app/adapters/presenter/tsv"
	"krohobor/app/usecases"
	"strings"
	"time"
)

type DumpCreate struct {
	UseCase usecases.DumpCreateInterface
}

func (d DumpCreate) Action(c *cli.Context, cfg config.AppConfig) error {
	dbs := c.String("dbname")

	var dbNames []string
	if dbs != "" {
		dbNames = strings.Split(dbs, ",")
	} else {
		dbs = "all"
	}

	filename := fmt.Sprintf("%s_%s.sql", dbs, time.Now().Format(cfg.Format))

	request := usecases.DumpCreateRequest{
		DbNames:  dbNames,
		Filename: filename,
	}

	resp, err := d.UseCase.Execute(request)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	pr := tsv.Presenter{}
	return pr.Print(resp)
}
