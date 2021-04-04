package archive

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Zip struct {
	Dir      string
	Password string
}

func NewZip(dir, password string) Zip {
	return Zip{dir, password}
}

func NewZipMock(dir, password string) Zip {
	if err := os.RemoveAll(dir); err != nil {
		panic(err)
	}

	mockDirPath := dir + "/mock"

	zip := Zip{dir, password}
	if err := os.MkdirAll(mockDirPath, os.ModePerm); err != nil {
		panic(err)
	}
	if err := zip.Archive(dir+"/test-mock-empty-folder.zip", mockDirPath); err != nil {
		panic(err)
	}

	d1 := []byte("hello")
	if err := ioutil.WriteFile(mockDirPath+"/file1", d1, 0644); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(mockDirPath+"/file2", d1, 0644); err != nil {
		panic(err)
	}

	if err := zip.Archive(dir+"/test-mock-with-password.zip", mockDirPath+"/file2"); err != nil {
		panic(err)
	}

	zipWithoutPass := Zip{dir, ""}
	if err := zipWithoutPass.Archive(dir+"/test-mock-without-password.zip", mockDirPath+"/file2"); err != nil {
		panic(err)
	}

	return zip
}

func (s Zip) Check() error {
	cmd := exec.Command("zip", "--version")
	_, err := cmd.Output()

	return err
}

func (s Zip) Archive(file, dir string) error {
	var cmd *exec.Cmd
	if "" != s.Password {
		cmd = exec.Command("zip", "-r", "--password", s.Password, file, dir)
	} else {
		cmd = exec.Command("zip", "-r", file, dir)
	}
	out, err := cmd.Output()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return errors.New(string(out))
		}
		return err
	}

	return nil
}

func (s Zip) Unarchive(file string) (string, error) {
	var cmd *exec.Cmd
	if "" != s.Password {
		cmd = exec.Command("unzip", "-j", "-o", "-P", s.Password, file, "-d", s.Dir)
	} else {
		cmd = exec.Command("unzip", "-j", "-o", file, "-d", s.Dir)
	}

	out, err := cmd.Output()
	if err != nil {
		if execErr, ok := err.(*exec.ExitError); ok {
			return "", errors.New(string(execErr.Stderr))
		}
		return "", err
	}

	return s.extract(out)
}

func (s Zip) extract(out []byte) (string, error) {
	findExt := "extracting: "
	find := findExt
	index := strings.Index(string(out), findExt)
	if index == -1 {
		findInf := "inflating: "
		index = strings.Index(string(out), findInf)
		if index == -1 {
			return "", errors.New("Not found inflating or extracting.")
		}
		find = findInf
	}
	res := strings.TrimSpace(string(out[index+len(find):]))
	return string(res), nil
}

func (s Zip) Ext() string {
	return ".zip"
}
