package archive

type Interface interface {
	Check() error
	Archive(string, string) error
	Unarchive(string) (string, error)
	Ext() string
}
