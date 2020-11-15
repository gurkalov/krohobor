package dump

import (
	"krohobor/app/adapters/database"
	"krohobor/app/adapters/storage"
)

type Dump struct {
	db database.Interface
	storage storage.Interface
}

func NewDump(db database.Interface, storage storage.Interface) Dump {
	return Dump{db, storage}
}

func (d Dump) Create(dbname, filename string) error {
	if err := d.db.Dump(dbname, d.storage.Filename(filename)); err != nil {
		return err
	}

	if err := d.storage.Write(filename); err != nil {
		return err
	}

	return nil
}

func (d Dump) CreateAll(string) error {
	panic("implement me")
}

func (d Dump) Restore(string) error {
	panic("implement me")
}

func (d Dump) Delete(string) error {
	panic("implement me")
}
