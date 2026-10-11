package usecases

import (
	"krohobor/app/adapters/database"
	"krohobor/app/adapters/storage"
	"krohobor/app/domain"
	"os"
)

type DumpRestoreInterface interface {
	Execute(DumpRestoreRequest) (DumpRestoreResponse, error)
}

type DumpRestore struct {
	db    database.Interface
	store storage.Interface

	// OnDownloaded is invoked after the dump is downloaded from storage and
	// before it is restored to the database, with the downloaded size in bytes.
	OnDownloaded func(size int)
}

type DumpRestoreRequest struct {
	Name     string
	DB       string
	Filename string
}

type DumpRestoreResponse struct {
	List []domain.Table
}

func NewDumpRestore(db database.Interface, store storage.Interface) *DumpRestore {
	return &DumpRestore{db: db, store: store}
}

func (dl *DumpRestore) Execute(request DumpRestoreRequest) (DumpRestoreResponse, error) {
	response := DumpRestoreResponse{}

	filename, err := dl.store.Read(request.Filename)
	if err != nil {
		return response, err
	}

	if dl.OnDownloaded != nil {
		size := 0
		if info, statErr := os.Stat(filename); statErr == nil {
			size = int(info.Size())
		}
		dl.OnDownloaded(size)
	}

	if err := dl.db.Restore(filename, request.Name); err != nil {
		return response, err
	}

	if err := dl.store.Clean(filename); err != nil {
		return response, err
	}

	return response, nil
}
