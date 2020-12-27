package usecases

import (
	"krohobor/app/adapters/database"
	"krohobor/app/domain"
)

type DbListInterface interface {
	Execute(DbListRequest) (DbListResponse, error)
}

type DbList struct {
	db database.Interface
}

type DbListRequest struct {}

type DbListResponse struct {
	List []domain.Database
}

func NewDbList(db database.Interface) *DbList {
	return &DbList{db}
}

func (dl *DbList) Execute(request DbListRequest) (DbListResponse, error) {
	response := DbListResponse{}

	resp, err := dl.db.List()
	if err != nil {
		return response, err
	}

	response.List = resp
	return response, nil
}
