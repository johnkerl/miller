package lib

import (
	"strconv"
)

func Itoa64(integer int64) string {
	return strconv.FormatInt(integer, 10)
}
