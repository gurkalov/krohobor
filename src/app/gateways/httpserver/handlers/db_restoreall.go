package handlers

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"krohobor/app/usecases"
	"net/http"
)

type DbRestoreAll struct {
	UseCase usecases.DbRestoreAllInterface
}

func (h DbRestoreAll) Handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		filename := "/tmp/download_backup/all.sql"
		archfile := "/tmp/download_backup.zip"
		archdir := "/tmp/download_backup"

		request := usecases.DbRestoreAllRequest{
			Backupname: ps.ByName("name"),
			Filename: filename,
			Archfile: archfile,
			Archdir: archdir,
		}

		resp, err := h.UseCase.Execute(request)
		if err != nil {
			_, _ = fmt.Fprint(w, err.Error())
			return
		}

		_, _ = fmt.Fprint(w, resp)
	}
}