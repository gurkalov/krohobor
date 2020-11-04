package usecases

import (
	"krohobor/app/adapters/archive"
	"krohobor/app/adapters/database"
	"krohobor/app/adapters/storage"
	"krohobor/app/domain"
)

type DbRestoreInterface interface {
	Execute(DbRestoreRequest) (DbRestoreResponse, error)
}

type DbRestore struct {
	db database.Interface
	arch  archive.Interface
	store storage.Interface
}

type DbRestoreRequest struct {
	Name string
	DB   string
	Filename string
	Archfile  string
	Archdir  string
}

type DbRestoreResponse struct {
	List []domain.Table
}

func NewDbRestore(db database.Interface, arch archive.Interface, store storage.Interface) *DbRestore {
	return &DbRestore{db, arch, store}
}

func (dl *DbRestore) Execute (request DbRestoreRequest) (DbRestoreResponse, error) {
	response := DbRestoreResponse{}

	if err := dl.store.Read(request.Name); err != nil {
		return response, err
	}

	if err := dl.arch.Unarchive(request.Archfile, request.Archdir); err != nil {
		return response, err
	}

	if err := dl.db.Restore(request.Filename); err != nil {
		return response, err
	}

	return response, nil
}
