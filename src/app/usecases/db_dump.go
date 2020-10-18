package usecases

import (
	"krohobor/app/adapters/archive"
	"krohobor/app/adapters/database"
	"krohobor/app/adapters/storage"
	"krohobor/app/domain"
)

type DbDumpInterface interface {
	Execute(DbDumpRequest) (DbDumpResponse, error)
}

type DbDump struct {
	db    database.Interface
	arch  archive.Interface
	store storage.Interface
}

type DbDumpRequest struct {
	Name string
	Filename string
}

type DbDumpResponse struct {
	List []domain.Table
}

func NewDbDump(db database.Interface, arch archive.Interface, store storage.Interface) *DbDump {
	return &DbDump{db, arch, store}
}

func (dd *DbDump) Execute (request DbDumpRequest) (DbDumpResponse, error) {
	response := DbDumpResponse{}

	if err := dd.db.Dump(request.Name, request.Filename); err != nil {
		return response, err
	}

	return response, nil
}
