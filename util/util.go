package util

import (
	"log"
	"strings"
)

func FisrtLower(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

func Sgo(fn func() error, errMsg string) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(errMsg, r)
			}
		}()
		err := fn()
		if err != nil {
			log.Println(errMsg, err)
		}
	}()
}
