package usecases

import (
	"krohobor/app/adapters/database"
)

type DbDeleteInterface interface {
	Execute(DbDeleteRequest) (DbDeleteResponse, error)
}

type DbDelete struct {
	db database.Interface
}

type DbDeleteRequest struct {
	Name  string
	Force bool
}

type DbDeleteResponse struct{}

func NewDbDelete(db database.Interface) *DbDelete {
	return &DbDelete{db}
}

func (dl *DbDelete) Execute(request DbDeleteRequest) (DbDeleteResponse, error) {
	response := DbDeleteResponse{}

	if err := dl.db.Drop(request.Name, request.Force); err != nil {
		return response, err
	}

	return response, nil
}
