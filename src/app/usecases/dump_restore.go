package usecases

import (
	"krohobor/app/adapters/database"
	"krohobor/app/adapters/storage"
	"krohobor/app/domain"
)

type DumpRestoreInterface interface {
	Execute(DumpRestoreRequest) (DumpRestoreResponse, error)
}

type DumpRestore struct {
	db database.Interface
	store storage.Interface
}

type DumpRestoreRequest struct {
	Name string
	DB   string
	Filename string
}

type DumpRestoreResponse struct {
	List []domain.Table
}

func NewDumpRestore(db database.Interface, store storage.Interface) *DumpRestore {
	return &DumpRestore{db, store}
}

func (dl *DumpRestore) Execute(request DumpRestoreRequest) (DumpRestoreResponse, error) {
	response := DumpRestoreResponse{}

	filename, err := dl.store.Read(request.Filename)
	if err != nil {
		return response, err
	}

	if err := dl.db.Restore(filename, request.Name); err != nil {
		return response, err
	}

	return response, nil
}
