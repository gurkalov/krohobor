package usecases

import (
	"krohobor/app/adapters/archive"
	"krohobor/app/adapters/database"
	"krohobor/app/adapters/storage"
	"krohobor/app/domain"
	"os"
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
	Dirname  string
	Archname string
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

	if err := dd.arch.Archive(request.Archname, request.Dirname); err != nil {
		return response, err
	}

	if err := dd.store.Write(request.Archname); err != nil {
		return response, err
	}

	if err := os.Remove(request.Filename); err != nil {
		return response, err
	}

	return response, nil
}
