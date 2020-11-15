package storage

import (
	"io/ioutil"
	"krohobor/app/adapters/archive"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	dir string
	archive archive.Interface
}

func NewFile(dir string, arch archive.Interface) File {
	return File{dir, arch}
}

func NewFileMock(dir string, arch archive.Interface) File {
	if err := os.RemoveAll(dir); err != nil {
		panic(err)
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}

	d1 := []byte("hello")
	if err := ioutil.WriteFile(dir + "/file1.txt", d1, 0644); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(dir + "/file2.txt", d1, 0644); err != nil {
		panic(err)
	}

	if arch != nil {
		archExt := arch.Ext()
		if err := arch.Archive(dir+"/file1.txt"+archExt, dir+"/file1.txt"); err != nil {
			panic(err)
		}

		if err := arch.Archive(dir+"/file2.txt"+archExt, dir+"/file2.txt"); err != nil {
			panic(err)
		}
	}

	return File{dir, arch}
}

func (s File) Check() error {
	if _, err := os.Stat(s.dir); os.IsNotExist(err) {
		return err
	}

	if s.archive != nil {
		return s.archive.Check()
	}

	return nil
}

func (s File) Filename(filename string) string {
	return s.dir + "/" + filename
}

func (s File) Read(filename string) (string, error) {
	archFile := filename
	if s.archive != nil {
		archFile = filename + s.archive.Ext()
	}

	if _, err := os.Stat(s.dir + "/" + archFile); os.IsNotExist(err) {
		return "", err
	}

	if s.archive != nil {
		return s.archive.Unarchive(s.dir + "/" + archFile)
	}

	return s.dir + "/" + filename, nil
}

func (s File) Write(filename string) error {
	if s.archive != nil {
		archFile := filename + s.archive.Ext()
		return s.archive.Archive(s.dir + "/" + archFile, s.dir + "/" + filename)
	}

	if _, err := os.Stat(s.dir + "/" + filename); err != nil {
		return err
	}

	return nil
}

func (s File) Delete(filename string) error {
	return os.RemoveAll(s.dir + "/" + filename)
}

func (s File) List() ([]string, error) {
	var result []string

	err := filepath.Walk(s.dir, func(path string, info os.FileInfo, err error) error {
		if s.dir != path {
			filename := path[len(s.dir)+1:]
			if s.archive != nil && s.archive.Ext() != "" {
				if s.archive.Ext() == filepath.Ext(filename) {
					filename = strings.TrimSuffix(filename, s.archive.Ext())
					result = append(result, filename)
				}
			} else {
				result = append(result, filename)
			}
		}
		return nil
	})
	if err != nil {
		return []string{}, err
	}

	return result, nil
}
