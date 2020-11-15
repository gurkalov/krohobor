package dump

type Interface interface {
	Create(string, string) error
	CreateAll(string) error
	Restore(string) error
	Delete(string) error
}
