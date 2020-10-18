package handlers

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"krohobor/app/usecases"
	"net/http"
)

type DbList struct {
	UseCase usecases.DbListInterface
}

func (h DbList) Handle() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		request := usecases.DbListRequest{}

		resp, err := h.UseCase.Execute(request)
		if err != nil {
			_, _ = fmt.Fprint(w, err.Error())
			return
		}

		_, _ = fmt.Fprint(w, resp)
	}
}
