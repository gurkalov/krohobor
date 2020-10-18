package handlers

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"krohobor/app/usecases"
	"net/http"
)

type DbRestore struct {
	UseCase usecases.DbRestoreInterface
}

func (h DbRestore) Handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		filename := "/tmp/backup/" + ps.ByName("db") + ".sql"

		request := usecases.DbRestoreRequest{
			Name: ps.ByName("db"),
			Filename: filename,
		}

		resp, err := h.UseCase.Execute(request)
		if err != nil {
			_, _ = fmt.Fprint(w, err.Error())
			return
		}

		_, _ = fmt.Fprint(w, resp)
	}
}