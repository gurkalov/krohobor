package usecases

import (
	"krohobor/app/adapters/database"
	"krohobor/app/adapters/storage"
)

type DumpCreateInterface interface {
	Execute(DumpCreateRequest) (DumpCreateResponse, error)
}

type DumpCreate struct {
	db    database.Interface
	store storage.Interface
}

type DumpCreateRequest struct {
	DbNames []string
	Filename string
}

type DumpCreateResponse struct {
	//List []domain.Table
}

func NewDumpCreate(db database.Interface, store storage.Interface) *DumpCreate {
	return &DumpCreate{db, store}
}

func (dd *DumpCreate) Execute(request DumpCreateRequest) (DumpCreateResponse, error) {
	response := DumpCreateResponse{}

	filename := dd.store.Filename(request.Filename)

	if len(request.DbNames) == 0 {
		if err := dd.db.DumpAll(filename); err != nil {
			return response, err
		}
	} else {
		if err := dd.db.Dump(request.DbNames[0], filename); err != nil {
			return response, err
		}
	}

	if err := dd.store.Write(request.Filename); err != nil {
		return response, err
	}

	return response, nil
}
