package database

import (
	"io/ioutil"
	"krohobor/app/domain"
	"os"
)

type Memory struct {
	list []domain.Database
}

func NewMemory() Memory {
	list := []domain.Database{
		{
			Name: "test1",
			Size: 100000,
		},
		{
			Name: "test2",
			Size: 100000,
		},
		{
			Name: "test3",
			Size: 100000,
		},
	}

	return Memory{list}
}

func (m Memory) Check() error {
	return nil
}

func (m Memory) List() ([]domain.Database, error) {
	return m.list, nil
}

func (m Memory) CreateDb(dbname string) error {
	db := domain.Database{
		Name: dbname,
		Size: 10000,
	}

	m.list = append(m.list, db)
	return nil
}

func (m Memory) Dump(dbname, filename string) error {
	sql := []byte("CREATE DATABASE " + dbname + ";")
	if err := ioutil.WriteFile(filename, sql, 0644); err != nil {
		return err
	}

	return nil
}

func (m Memory) DumpAll(filename string) error {
	sql := []byte("")
	for _, v := range m.list {
		sql = append(sql, []byte("CREATE DATABASE "+v.Name+";")...)
	}

	if err := ioutil.WriteFile(filename, sql, 0644); err != nil {
		return err
	}

	return nil
}

func (m Memory) Restore(filename, dbname string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return err
	}

	return nil
}

func (m Memory) Drop(dbname string) error {
	var newList []domain.Database
	for _, v := range m.list {
		if v.Name != dbname {
			newList = append(newList, v)
		}
	}

	m.list = newList
	return nil
}

func (m Memory) Tables(dbname string) ([]domain.Table, error) {
	list := []domain.Table{
		{
			Name:  "table1",
			Size:  10000,
			Count: 0,
		},
	}

	return list, nil
}

func (m Memory) Count(dbname, table string) (int, error) {

	return 0, nil
}
