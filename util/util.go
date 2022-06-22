package util

import (
	"strings"

	"github.com/kpango/glg"
)

func FisrtLower(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

func Sgo(fn func() error, errMsg string) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				glg.Error(errMsg, r)
			}
		}()
		err := fn()
		if err != nil {
			glg.Error(errMsg, err)
		}
	}()
}
