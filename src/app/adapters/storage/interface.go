package storage

import (
	"krohobor/app/domain"
)

type Interface interface {
	Check() error
	Filename(string) string
	Write(string) error
	Read(string) (string, error)
	Delete(string) error
	Clean(string) error
	List() ([]domain.Dump, error)
}
