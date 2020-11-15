package usecases

import (
	"krohobor/app/adapters/storage"
)

type BackupListInterface interface {
	Execute(BackupListRequest) (BackupListResponse, error)
}

type BackupList struct {
	store storage.Interface
}

type BackupListRequest struct {

}

type BackupListResponse struct {
	List []string
}

func NewBackupList(store storage.Interface) *BackupList {
	return &BackupList{store}
}

func (dl *BackupList) Execute(request BackupListRequest) (BackupListResponse, error) {
	response := BackupListResponse{}

	resp, err := dl.store.List()
	if err != nil {
		return response, err
	}

	response.List = resp
	return response, nil
}
