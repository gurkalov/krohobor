package archive

type Archive interface {
	Archive(string, string) error
	Unarchive(string, string) error
}
