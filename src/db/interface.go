package db

type Database interface {
	Init(string) error
	Count(string) int
	Create() error
	Drop() error
	List() []string
	Tables() []string
	Dump(string) error
	DumpAll(string) error
	Restore(string) error
	RestoreAll(string) error
	Info() map[string]int
}
