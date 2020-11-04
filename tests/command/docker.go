package command

import (
	"os/exec"
	"strings"
)

func Run(args ...string) (string, error) {
	baseArgs := []string{"exec", "-T", "app", "./krohobor"}
	baseArgs = append(baseArgs, args...)

	cmd := exec.Command("docker-compose", baseArgs...)
	cmd.Dir = "../"

	out, err := cmd.Output()

	return strings.TrimSuffix(string(out), "\n"), err
}
