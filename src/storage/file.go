package storage

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type File struct {
	Catalog string
}

func (s File) Read(filename string) error {
	fileParam := "/tmp/download_backup.zip"
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileParam, input, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s File) Write(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	dt := time.Now()
	nowDate := dt.Format("2006-01-02_15-04")
	key := "backup_" + nowDate + ".zip"

	os.MkdirAll(s.Catalog, 0755)
	err = ioutil.WriteFile(s.Catalog + "/" + key, file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s File) Delete(filename string) error {
	return os.RemoveAll(s.Catalog + "/" + filename)
}

func (s File) List() ([]string, error) {
	var result []string

	err := filepath.Walk(s.Catalog, func(path string, info os.FileInfo, err error) error {
		if s.Catalog != path {
			result = append(result, path)
		}
		return nil
	})
	if err != nil {
		return []string{}, err
	}

	return result, nil
}
