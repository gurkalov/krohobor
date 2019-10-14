package db

type Database interface {
	Init(string) error
	Create() error
	Drop() error
	List() []string
	Dump(string) error
	Restore(string) error
}
