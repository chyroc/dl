package helper

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func Rename(old, new string) error {
	_ = os.MkdirAll(filepath.Dir(new), 0o777)
	return os.Rename(old, new)
}

func ExistsPath(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func BuildPathIfNotExist(path string) error {
	ok, err := ExistsPath(path)
	if !ok {
		return os.MkdirAll(path, 0o644)
	}
	return err
}

func TrimFileExt(s string) string {
	if !strings.Contains(s, ".") {
		return s
	}
	ext := filepath.Ext(s)
	return s[:len(s)-len(ext)]
}

func TrimInvalidFilePathChars(path string) string {
	path = strings.TrimSpace(path)
	re := regexp.MustCompile("[\\\\/:*?\"<>|]")
	return re.ReplaceAllString(path, "")
}
