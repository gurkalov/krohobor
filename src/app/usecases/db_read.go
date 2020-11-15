package usecases

import (
	"krohobor/app/adapters/database"
	"krohobor/app/domain"
)

type DbReadInterface interface {
	Execute(DbReadRequest) (DbReadResponse, error)
}

type DbRead struct {
	db database.Interface
}

type DbReadRequest struct {
	Name string
}

type DbReadResponse struct {
	List []domain.Table
}

func NewDbRead(db database.Interface) *DbRead {
	return &DbRead{db}
}

func (dl *DbRead) Execute(request DbReadRequest) (DbReadResponse, error) {
	response := DbReadResponse{}

	resp, err := dl.db.Tables(request.Name)
	if err != nil {
		return response, err
	}

	response.List = resp
	return response, nil
}
