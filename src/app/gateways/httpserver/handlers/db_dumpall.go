package handlers

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"krohobor/app/usecases"
	"net/http"
)

type DbDumpAll struct {
	UseCase usecases.DbDumpAllInterface
}

func (h DbDumpAll) Handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		filename := "/tmp/backup/all.sql"
		dirname := "/tmp/backup"
		archname := "/tmp/backup_all.zip"

		request := usecases.DbDumpAllRequest{
			Filename: filename,
			Dirname: dirname,
			Archname: archname,
		}

		resp, err := h.UseCase.Execute(request)
		if err != nil {
			_, _ = fmt.Fprint(w, err.Error())
			return
		}

		_, _ = fmt.Fprint(w, resp)
	}
}