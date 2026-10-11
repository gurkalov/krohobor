package usecases

import (
	"krohobor/app/adapters/database"
	"krohobor/app/adapters/storage"
	"os"
)

type DumpCreateInterface interface {
	Execute(DumpCreateRequest) (DumpCreateResponse, error)
}

type DumpCreate struct {
	db    database.Interface
	store storage.Interface

	// OnDumped is invoked after the dump file is written by the database and
	// before it is archived and uploaded to storage, with the dump size in bytes.
	OnDumped func(size int)
}

type DumpCreateRequest struct {
	DbNames  []string
	Filename string
}

type DumpCreateResponse struct {
	//List []domain.Table
}

func NewDumpCreate(db database.Interface, store storage.Interface) *DumpCreate {
	return &DumpCreate{db: db, store: store}
}

func (dd *DumpCreate) Execute(request DumpCreateRequest) (DumpCreateResponse, error) {
	response := DumpCreateResponse{}

	filename := dd.store.Filename(request.Filename)

	if len(request.DbNames) == 0 {
		if err := dd.db.DumpAll(filename); err != nil {
			return response, err
		}
	} else {
		if err := dd.db.Dump(request.DbNames[0], filename); err != nil {
			return response, err
		}
	}

	if dd.OnDumped != nil {
		if info, statErr := os.Stat(filename); statErr == nil {
			dd.OnDumped(int(info.Size()))
		}
	}

	if err := dd.store.Write(request.Filename); err != nil {
		return response, err
	}

	return response, nil
}
