package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
	"xdt.com/hm-diag/config"
	"xdt.com/hm-diag/util"
)

var (
	PrivateAccessPaths = []string{
		"/inner/api/v1/miner/onboarding/txn",
	}
)

var WHITE_LIST = [...]string{
	"/favicon.ico",
	"/web",
	"/api/v1/login",
	"/inner/state",
	"/api/v1/password",
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Hotspot-Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func PrivateAccessMiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.Config().PublicAccess == config.CONF_ON {
			c.Next()
			return
		}

		inPrivate := false
		for _, p := range PrivateAccessPaths {
			glg.Debug(p, c.Request.URL.Path)
			if p == c.Request.URL.Path {
				inPrivate = true
				break
			}
		}
		if !inPrivate {
			c.Next()
			return
		}

		canAccess := false
		rIp, ok := c.RemoteIP()
		if !ok {
			glg.Debugf("get remote ip for request: %s from %s", c.Request.URL, c.Request.RemoteAddr)
			canAccess = false
		} else if util.IsPrivateIp(rIp) {
			glg.Debugf("is private ip for request: %s ", c.Request.URL)
			canAccess = true
		}
		if canAccess {
			c.Next()
		} else {
			c.JSON(401, RespBody{Code: 401, Message: "can't access on public environment"})
		}
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.Config().PublicAccess == config.CONF_ON {
			c.Next()
			return
		}

		inPrivate := false
		for _, p := range PrivateAccessPaths {
			glg.Debug(p, c.Request.URL.Path)
			if p == c.Request.URL.Path {
				inPrivate = true
				break
			}
		}
		if !inPrivate {
			c.Next()
			return
		}

		//信任地址或同ip免密
		canAccess := false
		rIp, ok := c.RemoteIP()
		if !ok {
			glg.Debugf("get remote ip for request: %s from %s", c.Request.URL, c.Request.RemoteAddr)
			canAccess = false
		} else if util.IsPrivateIp(rIp) {
			glg.Debugf("is private ip for request: %s ", c.Request.URL)
			canAccess = true
		}
		if canAccess {
			c.Next()
			return
		}

		p := c.Request.URL.Path

		//白名单免密
		white := false
		for _, w := range WHITE_LIST {
			if strings.HasPrefix(p, w) {
				white = true
			}
		}
		if p == "/" {
			white = true
		}
		if white {
			c.Next()
			return
		}

		//判断tk
		tk := c.GetHeader("Hotspot-Authorization")
		if tk == "" {
			tk = c.Query("hotspot_tk")
		}
		err := ValidateToken(tk)
		if err != nil {
			c.JSON(401, RespBody{Code: 401, Message: err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}
