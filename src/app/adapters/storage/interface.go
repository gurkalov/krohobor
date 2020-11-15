package storage

type Interface interface {
	Check() error
	Filename(string) string
	Write(string) error
	Read(string) (string, error)
	Delete(string) error
	List() ([]string, error)
}
