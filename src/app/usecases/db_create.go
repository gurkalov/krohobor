package usecases

import (
	"krohobor/app/adapters/database"
)

type DbCreateInterface interface {
	Execute(DbCreateRequest) (DbCreateResponse, error)
}

type DbCreate struct {
	db database.Interface
}

type DbCreateRequest struct {
	Name string
}

type DbCreateResponse struct{}

func NewDbCreate(db database.Interface) *DbCreate {
	return &DbCreate{db}
}

func (dl *DbCreate) Execute(request DbCreateRequest) (DbCreateResponse, error) {
	response := DbCreateResponse{}

	if err := dl.db.CreateDb(request.Name); err != nil {
		return response, err
	}

	return response, nil
}
