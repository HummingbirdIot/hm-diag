package util

import "strings"

func FisrtLower(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}
