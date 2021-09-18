package helper

import (
	"os"
	"path/filepath"
)

func Rename(old, new string) error {
	_ = os.MkdirAll(filepath.Dir(new), 0o777)
	return os.Rename(old, new)
}
