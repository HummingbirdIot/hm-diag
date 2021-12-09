package main

import (
	"embed"
	"io/fs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"xdt.com/hm-diag/ctrl"
	"xdt.com/hm-diag/diag"
	"xdt.com/hm-diag/proxy"
	"xdt.com/hm-diag/regist"
)

//go:embed web/release/*
var emFS embed.FS

var diagTask *diag.Task
var register *regist.Register

type RespBody struct {
	Data    interface{} `json:"data"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
}

func RespOK(data interface{}) RespBody {
	return RespBody{Data: data, Code: 200, Message: "OK"}
}

func route(r *gin.Engine, _diagTask *diag.Task, _register *regist.Register) {
	diagTask = _diagTask
	register = _register
	// r.LoadHTMLGlob("tmpl/*")

	RouteStatic(r)
	RoutePage(r)
	RouteState(r)
	RouteCtrl(r)
	RouteConfigProxy(r)
}

func RouteStatic(r *gin.Engine) {
	d, _ := fs.Sub(emFS, "web/release")
	r.StaticFS("/web", http.FS(d))
}

func RoutePage(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/web")
	})
}

func RouteCtrl(r *gin.Engine) {
	r.POST("/api/v1/device/reboot", func(c *gin.Context) {
		c.JSON(200, RespBody{
			Code:    200,
			Message: "receive reboot request, to reboot",
		})
		go ctrl.RebootDevice()
	})

	r.POST("/api/v1/miner/resync", func(c *gin.Context) {
		err := ctrl.ResyncMiner()
		if err != nil {
			c.JSON(500, RespBody{
				Code:    500,
				Message: "receive resync miner request, but got error:" + err.Error(),
			})
		} else {
			c.JSON(200, RespBody{
				Code:    200,
				Message: "OK",
			})
		}
	})
}

func RouteState(r *gin.Engine) {
	r.GET("/state", stateHandler)
	r.GET("/api/v1/device/state", deviceInfoHandler)
	r.GET("/api/v1/miner/state", minerInfoHandler)
	r.GET("/registInfo", registInfoHandler)
}

func RouteConfigProxy(r *gin.Engine) {
	r.GET("/api/v1/config/proxy", proxyGetHandler)
	r.POST("/api/v1/config/proxy", proxySetHandler)
}

func proxyGetHandler(c *gin.Context) {
	item := c.Query("item")
	var err error
	var p *proxy.ProxyItem
	if item == "gitRepo" {
		p, err = proxy.RepoProxy(opt.GitRepoDir)
	} else if item == "gitRelease" {
		p, err = proxy.ReleaseFileProxy(opt.GitRepoDir)
	}
	if err != nil {
		c.JSON(500, RespBody{Code: 500, Message: err.Error()})
	} else {
		c.JSON(200, RespOK(p))
	}
}

func proxySetHandler(c *gin.Context) {
	item := c.Query("item")
	if item != "gitRepo" && item != "gitRelease" {
		c.JSON(400, RespBody{Code: 400, Message: "query param 'item' should be 'gitRepo' or 'gitRelease'"})
		return
	}
	var proxyItem proxy.ProxyItem
	err := c.BindJSON(&proxyItem)
	if err != nil {
		c.JSON(400, RespBody{Code: 400, Message: "wrong request body:" + err.Error()})
		return
	}
	if item == "gitRepo" {
		err = proxy.SetRepoMirrorProxy(opt.GitRepoDir, proxyItem)
	} else if item == "gitRelease" {
		err = proxy.SetReleaseFileProxy(opt.GitRepoDir, proxyItem)
	}
	if err != nil {
		c.JSON(500, RespBody{Code: 500, Message: err.Error()})
	} else {
		c.JSON(200, RespOK(nil))
	}
}

func stateHandler(c *gin.Context) {
	d := diagTask.Data()
	if d.Data != nil {
		d.Data["aNotice"] = `do not use this api path "/" to integrate, use api under path "api/"`
	}
	c.JSON(200, RespOK(d.Data))
}

func minerInfoHandler(c *gin.Context) {
	var d diag.TaskData
	c.Query("cache")
	if c := c.Query("cache"); c == "true" {
		d = diagTask.MinerInfo()
	} else {
		d = diag.TaskData{Data: diagTask.FetchMinerInfo(), FetchTime: time.Now()}
	}
	c.JSON(200, RespOK(d.Data))
}

func deviceInfoHandler(c *gin.Context) {
	var d diag.TaskData
	if c := c.Query("cache"); c == "true" {
		d = diagTask.DeviceInfo()
	} else {
		d = diag.TaskData{Data: diagTask.FetchDeviceInfo(), FetchTime: time.Now()}
	}
	c.JSON(200, RespOK(d.Data))
}

func registInfoHandler(c *gin.Context) {
	d, err := register.GetRegistInfo()
	if err != nil {
		c.JSON(500, RespBody{Code: 500, Message: "error: " + err.Error()})
	} else {
		c.JSON(200, RespOK(d))
	}
}
