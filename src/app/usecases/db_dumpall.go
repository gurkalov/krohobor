package usecases

import (
	"krohobor/app/adapters/database"
	"krohobor/app/adapters/storage"
	"krohobor/app/domain"
)

type DbDumpAllInterface interface {
	Execute(DbDumpAllRequest) (DbDumpAllResponse, error)
}

type DbDumpAll struct {
	db    database.Interface
	store storage.Interface
}

type DbDumpAllRequest struct {
	Filename string
}

type DbDumpAllResponse struct {
	List []domain.Table
}

func NewDbDumpAll(db database.Interface, store storage.Interface) *DbDumpAll {
	return &DbDumpAll{db, store}
}

func (ddl *DbDumpAll) Execute(request DbDumpAllRequest) (DbDumpAllResponse, error) {
	response := DbDumpAllResponse{}

	if err := ddl.db.DumpAll(ddl.store.Filename(request.Filename)); err != nil {
		return response, err
	}

	if err := ddl.store.Write(request.Filename); err != nil {
		return response, err
	}

	return response, nil
}
