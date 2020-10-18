package usecases

import (
	"krohobor/app/adapters/database"
	"krohobor/app/domain"
)

type DbRestoreInterface interface {
	Execute(DbRestoreRequest) (DbRestoreResponse, error)
}

type DbRestore struct {
	db database.Interface
}

type DbRestoreRequest struct {
	Name string
	Filename string
}

type DbRestoreResponse struct {
	List []domain.Table
}

func NewDbRestore(db database.Interface) *DbRestore {
	return &DbRestore{db}
}

func (dl *DbRestore) Execute (request DbRestoreRequest) (DbRestoreResponse, error) {
	response := DbRestoreResponse{}

	if err := dl.db.Restore(request.Name, request.Filename); err != nil {
		return response, err
	}

	return response, nil
}
