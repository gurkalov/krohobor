package database

import (
	"krohobor/app/domain"
)

type Interface interface {
	List() ([]domain.Database, error)
	Tables(string) ([]domain.Table, error)
	Count(string, string) (int, error)
	CreateDb(string) error
	Dump(string, string) error
	DumpAll(string) error
	Drop(string) error
	Restore(string) error
	//Info(string) (map[string]int, error)
}
