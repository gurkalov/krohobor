package database

import (
	"krohobor/app/domain"
)

type Interface interface {
	Check() error
	List() ([]domain.Database, error)
	Tables(table string) ([]domain.Table, error)
	Count(string, string) (int, error)
	CreateDb(string) error
	Dump(string, string) error
	DumpAll(string) error
	Drop(dbname, target string) error
	Restore(filename string) error
	//Info(string) (map[string]int, error)
}
