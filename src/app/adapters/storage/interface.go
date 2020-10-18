package storage

type Interface interface {
	Write(string) error
	Read(string) error
	Delete(string) error
	List() ([]string, error)
}
