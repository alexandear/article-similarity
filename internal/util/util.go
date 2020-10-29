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

// Strip removes all non-alphanumeric characters from s.
func Strip(s []byte) []byte {
	n := 0

	for _, b := range s {
		if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || unicode.IsSpace(rune(b)) {
			s[n] = b
			n++
		}
	}

	return s[:n]
}

// IndexOf get the index of the given value in the given string slice, or -1 if not found.
func IndexOf(slice []string, value string) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}

	return -1
}

// Contains check if a string slice contains a value.
func Contains(slice []string, value string) bool {
	return IndexOf(slice, value) != -1
}
