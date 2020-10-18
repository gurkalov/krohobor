package archive

type Interface interface {
	Archive(string, string) error
	Unarchive(string, string) error
}
