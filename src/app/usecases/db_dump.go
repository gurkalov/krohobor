package usecases

import (
	"krohobor/app/adapters/database"
	"krohobor/app/adapters/storage"
)

type DbDumpInterface interface {
	Execute(DbDumpRequest) (DbDumpResponse, error)
}

type DbDump struct {
	db    database.Interface
	store storage.Interface
}

type DbDumpRequest struct {
	Name string
	Filename string
}

type DbDumpResponse struct {
	//List []domain.Table
}

func NewDbDump(db database.Interface, store storage.Interface) *DbDump {
	return &DbDump{db, store}
}

func (dd *DbDump) Execute(request DbDumpRequest) (DbDumpResponse, error) {
	response := DbDumpResponse{}

	if err := dd.db.Dump(request.Name, dd.store.Filename(request.Filename)); err != nil {
		return response, err
	}

	if err := dd.store.Write(request.Filename); err != nil {
		return response, err
	}

	return response, nil
}
