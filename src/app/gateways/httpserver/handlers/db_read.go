package handlers

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"krohobor/app/usecases"
	"net/http"
)

type DbRead struct {
	UseCase usecases.DbReadInterface
}

func (h DbRead) Handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		request := usecases.DbReadRequest{
			Name: ps.ByName("db"),
		}

		resp, err := h.UseCase.Execute(request)
		if err != nil {
			_, _ = fmt.Fprint(w, err.Error())
			return
		}

		_, _ = fmt.Fprint(w, resp)
	}
}