package main

import (
	"krohobor/db"
	"krohobor/storage"
	"os"
	"reflect"
	"testing"
)

func tearDown() {
	os.RemoveAll(cfg.App.Catalog)
	os.RemoveAll("/tmp/backup")
	os.RemoveAll("/tmp/backup.zip")
	os.RemoveAll("/tmp/download_backup")
	os.RemoveAll("/tmp/download_backup.zip")
}

func TestMain(m *testing.M) {
	InitConfig()
	InitManager()
	manager.Storage = &storage.File{cfg.App.Catalog}

	exitcode := m.Run()
	os.Exit(exitcode)
}

func TestBackupSuccess(t *testing.T) {
	tearDown()

	if err := manager.Backup([]string{"test1"}); err != nil {
		t.Errorf("Error %v.", err.Error())
	}

	list, err := manager.Storage.List()
	if err != nil {
		t.Errorf("Error %v.", err.Error())
	}

	actual := len(list)
	expected := 1
	if actual != expected {
		t.Errorf("Error len expected %v, but actual %v", expected, actual)
	}
}

func TestBackupError(t *testing.T) {
	tearDown()
	error := db.Postgres{"error"}
	error.Drop()

	err := manager.Backup([]string{"error"})
	if err.Error() != "exit status 1" {
		t.Errorf("Error %v.", err.Error())
	}
}

func TestRestoreSuccess(t *testing.T) {
	tearDown()

	if err := manager.Backup([]string{"test1"}); err != nil {
		t.Errorf("Error %v.", err.Error())
	}
	list, _ := manager.Storage.List()

	if err := manager.Restore([]string{"test1"}, list[0]); err != nil {
		t.Errorf("Error %v.", err.Error())
	}
}

func TestRestoreNoFileError(t *testing.T) {
	tearDown()
	error := db.Postgres{"error"}
	error.Drop()

	err := manager.Restore([]string{"error"}, "/test")
	if err.Error() != "open /test: no such file or directory" {
		t.Errorf("Error %v.", err.Error())
	}
}

func TestRestoreNoDatabaseError(t *testing.T) {
	tearDown()
	error := db.Postgres{"error"}
	error.Drop()

	if err := manager.Backup([]string{"test1"}); err != nil {
		t.Errorf("Error %v.", err.Error())
	}
	list, _ := manager.Storage.List()

	err := manager.Restore([]string{"error"}, list[0])
	if err.Error() != "stat /tmp/download_backup/backup_error.sql: no such file or directory" {
		t.Errorf("Error %v.", err.Error())
	}
}

func TestInfo(t *testing.T) {
	tearDown()
	test1 := db.Postgres{"test1"}
	actual := test1.Info()

	expected := map[string]int{
		"account1": 0,
		"link": 2,
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Error, expected: %v, but actual: %v", expected, actual)
	}
}

func TestFlow(t *testing.T) {
	tearDown()
	dbname := "test1"
	DB := db.Postgres{dbname}
	actual := DB.Info()

	expected := map[string]int{
		"account1": 0,
		"link": 2,
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nExpected: %v \nActual:   %v", expected, actual)
	}

	if err := manager.Backup([]string{dbname}); err != nil {
		t.Errorf("Error %v.", err.Error())
	}

	//DB.Drop()
	//droppedActual := DB.Info()
	//droppedExpected := map[string]int{}
	//if !reflect.DeepEqual(droppedExpected, droppedActual) {
	//	t.Errorf("\nExpected: %v \nActual:   %v", droppedExpected, droppedActual)
	//}

	list, _ := manager.Storage.List()
	if err := manager.Restore([]string{dbname}, list[0]); err != nil {
		t.Errorf("Error %v.", err.Error())
	}

	restoreActual := DB.Info()
	restoreExpected := expected
	if !reflect.DeepEqual(restoreExpected, restoreActual) {
		t.Errorf("\nExpected: %v \nActual:   %v", restoreExpected, restoreActual)
	}
}

func TestBackupAllSuccess(t *testing.T) {
	tearDown()

	if err := manager.BackupAll(); err != nil {
		t.Errorf("Error %v.", err.Error())
	}

	list, err := manager.Storage.List()
	if err != nil {
		t.Errorf("Error %v.", err.Error())
	}

	actual := len(list)
	expected := 1
	if actual != expected {
		t.Errorf("Error len expected %v, but actual %v", expected, actual)
	}
}

func TestRestoreAllSuccess(t *testing.T) {
	tearDown()

	if err := manager.BackupAll(); err != nil {
		t.Errorf("Error %v.", err.Error())
	}
	list, _ := manager.Storage.List()

	if err := manager.RestoreAll(list[0]); err != nil {
		t.Errorf("Error %v.", err.Error())
	}
}
