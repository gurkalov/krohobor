package usecases

import (
	"krohobor/app/adapters/storage"
	"krohobor/app/domain"
)

type DumpDeleteInterface interface {
	Execute(DumpDeleteRequest) (DumpDeleteResponse, error)
}

type DumpDelete struct {
	store storage.Interface
}

type DumpDeleteRequest struct {
	Name string
}

type DumpDeleteResponse struct {
	List []domain.Table
}

func NewDumpDelete(store storage.Interface) *DumpDelete {
	return &DumpDelete{store}
}

func (dl *DumpDelete) Execute(request DumpDeleteRequest) (DumpDeleteResponse, error) {
	response := DumpDeleteResponse{}

	if err := dl.store.Delete(request.Name); err != nil {
		return response, err
	}

	return response, nil
}
