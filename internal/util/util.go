package util

import (
	"unicode"
)

// Min returns the value of the smallest argument.
func Min(args ...int) int {
	if len(args) == 0 {
		return 0
	}
	if len(args) == 1 {
		return args[0]
	}

	min := args[0]
	for _, arg := range args[1:] {
		if min > arg {
			min = arg
		}
	}

	return min
}

// Max returns the value of the largest argument.
func Max(args ...int) int {
	if len(args) == 0 {
		return 0
	}
	if len(args) == 1 {
		return args[0]
	}

	max := args[0]
	for _, arg := range args[1:] {
		if max < arg {
			max = arg
		}
	}

	return max
}

// Strip removes all non-alphanumeric characters from s
func Strip(s []byte) []byte {
	n := 0
	for _, b := range s {
		if ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z') || ('0' <= b && b <= '9') || unicode.IsSpace(rune(b)) {
			s[n] = b
			n++
		}
	}
	return s[:n]
}
