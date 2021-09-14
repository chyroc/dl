package helper

import (
	"strconv"
)

func MayStringToInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}
