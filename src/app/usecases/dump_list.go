package usecases

import (
	"krohobor/app/adapters/storage"
	"krohobor/app/domain"
)

type DumpListInterface interface {
	Execute(DumpListRequest) (DumpListResponse, error)
}

type DumpList struct {
	store storage.Interface
}

type DumpListRequest struct {
}

type DumpListResponse struct {
	List []domain.Dump
}

func NewDumpList(store storage.Interface) *DumpList {
	return &DumpList{store}
}

func (dl *DumpList) Execute(request DumpListRequest) (DumpListResponse, error) {
	response := DumpListResponse{}

	resp, err := dl.store.List()
	if err != nil {
		return response, err
	}

	response.List = resp
	return response, nil
}
