package database

import (
	"krohobor/app/domain"
)

type Interface interface {
	Check(target string) error
	List(target string) ([]domain.Database, error)
	Tables(table, target string) ([]domain.Table, error)
	Count(string, string) (int, error)
	CreateDb(string) error
	Dump(string, string) error
	DumpAll(string) error
	Drop(dbname, target string) error
	Restore(filename, target string) error
	//Info(string) (map[string]int, error)
}
