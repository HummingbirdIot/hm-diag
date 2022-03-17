package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"xdt.com/hm-diag/config"
	"xdt.com/hm-diag/util"
)

var (
	PrivateAccessPaths = []string{
		"/inner/api/v1/miner/onboarding/txn",
	}
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
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
			log.Println(p, c.Request.URL.Path)
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
			log.Printf("get remote ip for request: %s from %s", c.Request.URL, c.Request.RemoteAddr)
			canAccess = false
		} else if util.IsPrivateIp(rIp) {
			log.Printf("is private ip for request: %s ", c.Request.URL)
			canAccess = true
		}
		if canAccess {
			c.Next()
		} else {
			c.JSON(401, RespBody{Code: 401, Message: "can't access on public environment"})
		}
	}
}