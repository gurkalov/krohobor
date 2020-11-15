package usecases

import (
	"krohobor/app/adapters/database"
	"krohobor/app/adapters/storage"
	"krohobor/app/domain"
)

type DbRestoreInterface interface {
	Execute(DbRestoreRequest) (DbRestoreResponse, error)
}

type DbRestore struct {
	db database.Interface
	store storage.Interface
}

type DbRestoreRequest struct {
	Name string
	DB   string
	Filename string
}

type DbRestoreResponse struct {
	List []domain.Table
}

func NewDbRestore(db database.Interface, store storage.Interface) *DbRestore {
	return &DbRestore{db, store}
}

func (dl *DbRestore) Execute(request DbRestoreRequest) (DbRestoreResponse, error) {
	response := DbRestoreResponse{}

	filename, err := dl.store.Read(request.Filename)
	if err != nil {
		return response, err
	}

	if err := dl.db.Restore(filename); err != nil {
		return response, err
	}

	return response, nil
}
