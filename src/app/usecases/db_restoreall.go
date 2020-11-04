package usecases

import (
	"krohobor/app/adapters/archive"
	"krohobor/app/adapters/database"
	"krohobor/app/adapters/storage"
	"krohobor/app/domain"
)

type DbRestoreAllInterface interface {
	Execute(DbRestoreAllRequest) (DbRestoreAllResponse, error)
}

type DbRestoreAll struct {
	db    database.Interface
	arch  archive.Interface
	store storage.Interface
}

type DbRestoreAllRequest struct {
	Backupname string
	Filename string
	Archfile  string
	Archdir  string
}

type DbRestoreAllResponse struct {
	List []domain.Table
}

func NewDbRestoreAll(db database.Interface, arch archive.Interface, store storage.Interface) *DbRestoreAll {
	return &DbRestoreAll{db, arch, store}
}

func (d *DbRestoreAll) Execute (request DbRestoreAllRequest) (DbRestoreAllResponse, error) {
	response := DbRestoreAllResponse{}

	if err := d.store.Read(request.Backupname); err != nil {
		return response, err
	}

	if err := d.arch.Unarchive(request.Archfile, request.Archdir); err != nil {
		return response, err
	}

	//if err := d.db.RestoreAll(request.Filename); err != nil {
	//	return response, err
	//}

	return response, nil
}
