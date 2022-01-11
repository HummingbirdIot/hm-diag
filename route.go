package main

import (
	"embed"
	"encoding/base64"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"xdt.com/hm-diag/ctrl"
	"xdt.com/hm-diag/devdis"
	"xdt.com/hm-diag/diag"
	"xdt.com/hm-diag/diag/miner"
	"xdt.com/hm-diag/regist"
	"xdt.com/hm-diag/util"
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
	// r.OPTIONS("/", CORSMiddleware())
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

	r.POST("/api/v1/miner/restart", func(c *gin.Context) {
		err := ctrl.RestartMiner()
		if err != nil {
			c.JSON(500, RespBody{
				Code:    500,
				Message: "receive restart miner request, but got error:" + err.Error(),
			})
		} else {
			c.JSON(200, RespBody{
				Code:    200,
				Message: "OK",
			})
		}
	})

	r.POST("/api/v1/miner/snapshot", func(c *gin.Context) {
		log.Println("to snapshot")
		ctrl.SnapshotTake()
		c.JSON(200, RespOK(nil))
	})

	r.GET("/api/v1/miner/snapshot/state", func(c *gin.Context) {
		log.Println("to get snapshot state")
		res, err := ctrl.SnapshotState()
		if err != nil {
			log.Printf("get snapshot state error: %+v\n", err)
			c.JSON(500, RespBody{Code: 500, Message: err.Error()})
		} else {
			c.JSON(200, RespOK(res))
		}
	})

	r.GET("/api/v1/miner/snapshot/file/:name", func(c *gin.Context) {
		log.Println("to get snapshot file")
		f := c.Param("name")
		b, err := base64.StdEncoding.DecodeString(f)
		if err != nil {
			log.Println("get snapshot, wrong path param", f)
			c.JSON(400, RespBody{Code: 400, Message: "wrong path param"})
		}
		f = string(b)
		if !strings.HasPrefix(f, "/tmp/") {
			log.Println("get snapshot, wrong path param", f)
			c.JSON(400, RespBody{Code: 400, Message: "wrong path param"})
		}
		c.Header("Content-Disposition", "attachment; filename=\""+path.Base(f)+"\"")
		c.File(f)
	})

	r.POST("/api/v1/miner/snapshot/apply", func(c *gin.Context) {
		log.Println("to apply snapshot")
		f, err := c.FormFile("file")
		if err != nil {
			log.Println("get form file error", err)
			c.JSON(400, RespBody{
				Code:    400,
				Message: err.Error(),
			})
			return
		}
		tmpF := "/tmp/" + strconv.FormatInt(time.Now().UnixNano(), 10)
		c.SaveUploadedFile(f, tmpF)
		ctrl.SnapshotLoad(tmpF)
		c.JSON(200, RespOK(nil))
	})

	r.POST("/api/v1/docker/reset", func(c *gin.Context) {
		ctrl.DockerReset()
		c.JSON(200, RespOK(nil))
	})

	r.POST("/api/v1/workspace/reset", func(c *gin.Context) {
		err := ctrl.GitRepoReset()
		if err != nil {
			log.Println(err)
			c.JSON(500, RespBody{Code: 500, Message: err.Error()})
		} else {
			c.JSON(200, RespOK(nil))
		}
	})
	r.GET("/api/v1/workspace/update", func(c *gin.Context) {
		r, err := ctrl.IsGitRepoToUpdate()
		if err != nil {
			c.JSON(500, RespBody{Code: 500, Message: err.Error()})
		} else {
			c.JSON(200, RespOK(r))
		}
	})
	r.POST("/api/v1/workspace/update", func(c *gin.Context) {
		ctrl.ExecMainUpdate()
		c.JSON(200, RespOK(nil))
	})
}

func RouteState(r *gin.Engine) {
	r.GET("/state", stateHandler)
	r.GET("/api/v1/device/state", deviceInfoHandler)
	r.GET("/api/v1/miner/state", minerInfoHandler)
	r.GET("/registInfo", registInfoHandler)
	r.GET("/api/v1beta/miner/log", minerLogHandler)
	r.GET("/api/v1/lan/hotspot", func(c *gin.Context) {
		c.JSON(http.StatusOK, RespOK(devdis.Services()))
	})
	r.GET("/api/v1/proxy/heliumApi", heliumApiProxyHandler)
}

func RouteConfigProxy(r *gin.Engine) {
	r.GET("/api/v1/config/proxy", proxyGetHandler)
	r.POST("/api/v1/config/proxy", proxySetHandler)
}

func proxyGetHandler(c *gin.Context) {
	item := c.Query("item")
	var err error
	var p *ctrl.ProxyItem
	if item == "gitRepo" {
		p, err = ctrl.RepoProxy(opt.GitRepoDir)
	} else if item == "gitRelease" {
		p, err = ctrl.ReleaseFileProxy(opt.GitRepoDir)
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
	var proxyItem ctrl.ProxyItem
	err := c.BindJSON(&proxyItem)
	if err != nil {
		c.JSON(400, RespBody{Code: 400, Message: "wrong request body:" + err.Error()})
		return
	}
	if item == "gitRepo" {
		err = ctrl.SetRepoMirrorProxy(opt.GitRepoDir, proxyItem)
	} else if item == "gitRelease" {
		err = ctrl.SetReleaseFileProxy(opt.GitRepoDir, proxyItem)
	}
	if err != nil {
		log.Println(err)
		c.JSON(500, RespBody{Code: 500, Message: err.Error()})
	} else {
		c.JSON(200, RespOK(nil))
	}
}

func stateHandler(c *gin.Context) {
	var res map[string]interface{}
	if c := c.Query("cache"); c == "true" {
		res = diagTask.Data().Data
		res["time"] = diagTask.Data().FetchTime
	} else {
		res = map[string]interface{}{
			"device": diagTask.FetchDeviceInfo(),
			"miner":  diagTask.FetchMinerInfo(),
		}
	}
	res["notice"] = `do not use this api path "/" to integrate, use api under path "api/"`
	c.JSON(200, RespOK(res))
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

func minerLogHandler(c *gin.Context) {
	var since time.Time = time.Now().Add(time.Minute * time.Duration(5))
	var until time.Time = time.Now()
	if st, err := time.Parse("2006-01-02T15:04:05", c.Query("since")); err == nil {
		since = st
	} else {
		log.Println("query miner log, convert since time error: ", err)
	}
	if tt, err := time.Parse("2006-01-02T15:04:05", c.Query("until")); err == nil {
		until = tt
	} else {
		log.Println("query miner log, convert until time error: ", err)
	}
	filter := c.Query("filter")
	log.Printf("query miner log, since: %s, until: %s, filter: %s", since, until, filter)
	l, err := miner.MinerLog(since, until, filter)
	if err == nil {
		c.JSON(200, RespOK(l))
	} else {
		c.JSON(500, RespBody{Code: 500, Message: err.Error()})
		log.Println("query miner log error: ", err)
	}
}

func heliumApiProxyHandler(c *gin.Context) {
	path := c.Query("path")
	s, err := util.HeliumApiProxy(path)
	if err != nil {
		log.Printf("helium proxy apt %s error: %s", path, err)
		c.JSON(500, RespBody{Code: 500, Message: err.Error()})
	} else {
		c.JSON(200, RespOK(s))
	}
}
