package api

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"
	"time"

	"xdt.com/hm-diag/config"
)

func GenToken() string {
	ts := time.Now().UnixMilli()
	tk := fmt.Sprintf("%s.%d", GenTokenSecret(ts), ts)
	return tk
}

func GenTokenSecret(ts int64) string {
	p := fmt.Sprintf("%s %d", config.Config().Password, ts)
	s := md5.Sum([]byte(p))

	return fmt.Sprintf("%x", s)
}

const TOKEN_EXPIRE_MS = 3600 * 24 * 30 * 1000

func ValidateToken(token string) error {
	a := strings.Split(token, ".")
	if len(a) != 2 {
		return fmt.Errorf("invalid token")
	}
	sum := a[0]
	ts, err := strconv.ParseInt(a[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid token")
	}
	if time.Now().UnixMilli()-ts > TOKEN_EXPIRE_MS {
		return fmt.Errorf("token expired")
	}
	if GenTokenSecret(ts) != sum {
		return fmt.Errorf("invalid token")
	}
	return nil
}
