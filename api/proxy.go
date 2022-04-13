package api

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"xdt.com/hm-diag/config"
	"xdt.com/hm-diag/ctrl"
	"xdt.com/hm-diag/util"
)

func gitProxyGet(c *gin.Context) {
	conf := config.Config()
	item := c.Query("item")
	var err error
	var p *ctrl.ProxyItem
	if item == "gitRepo" {
		p, err = ctrl.RepoProxy(conf.GitRepoDir)
	} else if item == "gitRelease" {
		p, err = ctrl.ReleaseFileProxy(conf.GitRepoDir)
	}
	if err != nil {
		fmt.Println(err)
		c.JSON(500, RespBody{Code: 500, Message: err.Error()})
	} else {
		c.JSON(200, RespOK(p))
	}
}

func gitProxySet(c *gin.Context) {
	conf := config.Config()
	item := c.Query("item")
	if item != "gitRepo" && item != "gitRelease" {
		c.JSON(400, RespBody{Code: 400, Message: "query param 'item' should be 'gitRepo' or 'gitRelease'"})
		return
	}
	var proxyItem ctrl.ProxyItem
	err := c.BindJSON(&proxyItem)
	if err != nil {
		c.JSON(400, RespBody{Code: 400, Message: "wrong request body:" + err.Error()})
		return
	}
	if item == "gitRepo" {
		err = ctrl.SetRepoMirrorProxy(conf.GitRepoDir, proxyItem)
	} else if item == "gitRelease" {
		err = ctrl.SetReleaseFileProxy(conf.GitRepoDir, proxyItem)
	}
	if err != nil {
		log.Println(err)
		c.JSON(500, RespBody{Code: 500, Message: err.Error()})
	} else {
		c.JSON(200, RespOK(nil))
	}
}

func heliumApiProxy(c *gin.Context) {
	path := c.Query("path")
	s, err := util.HeliumApiProxy(path)
	if err != nil {
		log.Printf("helium proxy apt %s error: %s", path, err)
		c.JSON(500, RespBody{Code: 500, Message: err.Error()})
	} else {
		c.JSON(200, RespOK(s))
	}
}
