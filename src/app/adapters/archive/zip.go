package archive

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type Zip struct {
	Dir string
	Password string
}

func NewZip(dir, password string) Zip {
	return Zip{dir,  password}
}

func NewZipMock(dir, password, mockDir string) Zip {
	mockDirPath := dir + "/" + mockDir

	if err := os.MkdirAll(mockDirPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	d1 := []byte("hello\ngo\n")
	if err := ioutil.WriteFile(mockDirPath + "/file1", d1, 0644); err != nil {
		log.Fatal(err)
	}

	return Zip{dir,  password}
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
	if err := os.RemoveAll(s.Dir); err != nil {
		return "", err
	}

	var cmd *exec.Cmd
	if "" != s.Password {
		cmd = exec.Command("unzip", "-j", "-P", s.Password, file, "-d", s.Dir)
	} else {
		cmd = exec.Command("unzip", "-j", file, "-d", s.Dir)
	}

	_, err := cmd.Output()
	if err != nil {
		if execErr, ok := err.(*exec.ExitError); ok {
			return "", errors.New(string(execErr.Stderr))
		}
		return "", err
	}

	return s.Dir + "file", nil
}
