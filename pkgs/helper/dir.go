package helper

import (
	"os"
	"path/filepath"
	"strings"
)

func ResolveDirOrCurrent(dir string) (string, error) {
	if strings.HasPrefix(dir, "/") {
		return strings.TrimRight(dir, "/"), nil
	}
	if strings.HasPrefix(dir, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return strings.TrimRight(filepath.Join(home, dir[2:]), "~"), nil
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(filepath.Join(currentDir, dir), "/"), nil
}
