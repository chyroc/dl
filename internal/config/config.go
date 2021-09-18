package config

import (
	"os"
)

var (
	Home    string
	WorkDir string
)

func isInTest() bool {
	return os.Getenv("is_in_test") != ""
}

func init() {
	h, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	Home = h

	WorkDir = Home + "/.dl"
	_ = os.MkdirAll(WorkDir, 0o777)
}
