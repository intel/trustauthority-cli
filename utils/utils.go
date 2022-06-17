package utils

import (
	"bufio"
	"github.com/pkg/errors"
	"os"
	"strings"
)

func ReadAnswerFileToEnv(filename string) error {
	fin, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err, "Failed to load answer file")
	}
	scanner := bufio.NewScanner(fin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" ||
			strings.HasPrefix(line, "#") {
			continue
		}
		equalSign := strings.Index(line, "=")
		if equalSign > 0 {
			key := line[0:equalSign]
			val := line[equalSign+1:]
			if key != "" &&
				val != "" {
				err = os.Setenv(key, val)
				if err != nil {
					return errors.Wrap(err, "Failed to set ENV")
				}
			}
		}
	}
	return nil
}
