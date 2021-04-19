package command

import "strings"

func ToStrings(val string) ([][]string, error) {
	result := make([][]string, 0)
	lines := strings.Split(val, "\n")
	for _, line := range lines {
		result = append(result, strings.Split(line, "\t"))
	}

	return result, nil
}
