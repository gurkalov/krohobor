package main

import (
	"krohobor/archive"
	"krohobor/db"
	"krohobor/storage"
	"os"
	"os/exec"
	"strings"
)

type Manager struct {
	password string
	Storage storage.Storage
	DB db.Database
	Archive archive.Archive
}

var manager Manager

func InitManager() {
	manager = Manager{cfg.App.Password,
		&storage.AwsS3{cfg.App.Catalog},
		&db.Postgres{},
		&archive.Zip{cfg.App.Password},
	}
}

func (s Manager) Backup(dbname []string) error {
	archFilename := "/tmp/backup.zip"
	archDirectory := "/tmp/backup"
	os.MkdirAll(archDirectory, 0755)
	for _, v := range dbname {
		db := s.DB
		db.Init(v)
		if err := db.Dump(archDirectory + "/backup_" + v + ".sql"); err != nil {
			return err
		}
	}

	if err := s.Archive.Archive(archFilename, archDirectory); err != nil {
		return err
	}

	if err := s.Storage.Write(archFilename); err != nil {
		return err
	}

	if err := os.Remove(archFilename); err != nil {
		return err
	}

	return nil
}

func (s Manager) BackupAll() error {
	archFilename := "/tmp/backup.zip"
	archDirectory := "/tmp/backup"
	os.MkdirAll(archDirectory, 0755)

	db := s.DB
	db.Init("postgres")
	if err := db.DumpAll(archDirectory + "/backup_all.sql"); err != nil {
		return err
	}

	if err := s.Archive.Archive(archFilename, archDirectory); err != nil {
		return err
	}

	if err := s.Storage.Write(archFilename); err != nil {
		return err
	}

	if err := os.Remove(archFilename); err != nil {
		return err
	}

	return nil
}

func (s Manager) Restore(dbname []string, filename string) error {
	if err := s.Storage.Read(filename); err != nil {
		return err
	}

	archDirectory := "/tmp/download_backup"
	if err := s.Archive.Unarchive("/tmp/download_backup.zip", archDirectory); err != nil {
		return err
	}

	for _, v := range dbname {
		db := s.DB
		db.Init(v)
		db.Drop()
		db.Create()
		if err := db.Restore(archDirectory + "/backup_" + v + ".sql"); err != nil {
			return err
		}
	}

	return nil
}

func (s Manager) RestoreAll(filename string) error {
	if err := s.Storage.Read(filename); err != nil {
		return err
	}

	archDirectory := "/tmp/download_backup"
	if err := s.Archive.Unarchive("/tmp/download_backup.zip", archDirectory); err != nil {
		return err
	}

	db := s.DB
	db.Init("postgres")
	if err := db.RestoreAll(archDirectory + "/backup_all.sql"); err != nil {
		return err
	}

	return nil
}

func (s Manager) List() []string {
	cmd := exec.Command("psql",
		"-t", "-A", `-F","`,
		"-c", "SELECT datname FROM pg_database WHERE datname NOT IN ('postgres', 'template0', 'template1', 'template2');")
	stdout, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	return strings.Split(strings.Trim(string(stdout), "\n"), "\n")
}
