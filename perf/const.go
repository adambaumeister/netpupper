package perf

import (
	"fmt"
	"strconv"
)

const GIGABYTE = 1000000000
const MEGABYTE = 1000000

/*
Convert a string with a byte delimiter to a byte len
	1K = 1000
	1M = 1000000
	etc..
*/
func StringToByte(s string) uint64 {
	sl := len(s)
	switch string(s[sl-1]) {
	case "M":
		v, _ := strconv.Atoi(s[:sl-1])
		return uint64(v * MEGABYTE)
	case "G":
		v, _ := strconv.Atoi(s[:sl-1])
		return uint64(v * GIGABYTE)
	default:
		return 1
	}
}

func ByteToString(b uint64) string {
	switch {
	case b > GIGABYTE:
		return fmt.Sprintf("%vG", b/GIGABYTE)
	case b > MEGABYTE:
		return fmt.Sprintf("%vM", b/MEGABYTE)
	default:
		return string(b)
	}
}
