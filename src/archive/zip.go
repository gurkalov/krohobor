package archive

import (
	"os"
	"os/exec"
)

type Zip struct {
	Password string
}

func (s Zip) Archive(file, dir string) error {
	var cmd *exec.Cmd
	if "" != s.Password {
		cmd = exec.Command("zip", "-r", "--password", s.Password, file, dir)
	} else {
		cmd = exec.Command("zip", "-r", file, dir)
	}
	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

func (s Zip) Unarchive(file, dir string) error {
	os.RemoveAll(dir)

	var cmd *exec.Cmd
	if "" != s.Password {
		cmd = exec.Command("unzip", "-j", "-P", s.Password, file, "-d", dir)
	} else {
		cmd = exec.Command("unzip", "-j", file, "-d", dir)
	}

	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}
