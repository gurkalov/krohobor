package usecases

import (
	"krohobor/app/adapters/archive"
	"krohobor/app/adapters/database"
	"krohobor/app/adapters/storage"
	"krohobor/app/domain"
	"os"
)

type DbDumpAllInterface interface {
	Execute(DbDumpAllRequest) (DbDumpAllResponse, error)
}

type DbDumpAll struct {
	db    database.Interface
	arch  archive.Interface
	store storage.Interface
}

type DbDumpAllRequest struct {
	Filename string
	Dirname  string
	Archname string
}

type DbDumpAllResponse struct {
	List []domain.Table
}

func NewDbDumpAll(db database.Interface, arch archive.Interface, store storage.Interface) *DbDumpAll {
	return &DbDumpAll{db, arch, store}
}

func (ddl *DbDumpAll) Execute (request DbDumpAllRequest) (DbDumpAllResponse, error) {
	response := DbDumpAllResponse{}

	if err := ddl.db.DumpAll(request.Filename); err != nil {
		return response, err
	}

	if err := ddl.arch.Archive(request.Archname, request.Dirname); err != nil {
		return response, err
	}

	if err := ddl.store.Write(request.Archname); err != nil {
		return response, err
	}

	if err := os.Remove(request.Filename); err != nil {
		return response, err
	}

	return response, nil
}
