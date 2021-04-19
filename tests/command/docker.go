package command

import (
	"os/exec"
	"strings"
	"time"
)

func Run(args ...string) (string, error) {
	baseArgs := []string{"exec", "-T", "app", "./krohobor"}
	baseArgs = append(baseArgs, args...)

	cmd := exec.Command("docker-compose", baseArgs...)
	cmd.Dir = "../"

	out, err := cmd.Output()

	return strings.TrimSuffix(string(out), "\n"), err
}

func TearDown() error {
	cmdDown := exec.Command("docker-compose", "down", "-v")
	cmdDown.Dir = "../"

	_, err := cmdDown.Output()
	if err != nil {
		//do nothing
	}

	cmdUp := exec.Command("docker-compose", "up", "-d")
	cmdUp.Dir = "../"

	_, err = cmdUp.Output()
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	return nil
}
