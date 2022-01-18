package api

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"xdt.com/hm-diag/diag"
	"xdt.com/hm-diag/regist"
)

type RespBody struct {
	Data    interface{} `json:"data"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
}

func RespOK(data interface{}) RespBody {
	return RespBody{Data: data, Code: 200, Message: "OK"}
}

var diagTask *diag.Task
var register *regist.Register
var webFS embed.FS

func Route(r *gin.Engine, webFiles embed.FS) {
	webFS = webFiles
	diagTask = diag.TaskInstance()
	register = regist.Instance()

	// web static files
	d, _ := fs.Sub(webFS, "web/release")
	r.StaticFS("/web", http.FS(d))
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/web")
	})

	// hotspot
	r.GET("/api/v1/device/state", deviceInfo)
	r.POST("/api/v1/device/reboot", rebootDevice)
	r.GET("/api/v1/lan/hotspot", lanHotspotsLoad)

	// miner
	r.POST("/api/v1/miner/resync", resyncMiner)
	r.POST("/api/v1/miner/restart", restartMiner)
	r.POST("/api/v1/miner/snapshot", snapshotTake)
	r.GET("/api/v1/miner/snapshot/state", snapshotState)
	r.GET("/api/v1/miner/snapshot/file/:name", snapshotDownload)
	r.POST("/api/v1/miner/snapshot/apply", snapshotApply)
	r.GET("/api/v1/miner/state", minerInfo)
	r.GET("/api/v1beta/miner/log", minerLogQuery)

	// workspace
	r.POST("/api/v1/workspace/reset", workspaceReset)
	r.GET("/api/v1/workspace/update", workspaceUpdateAvailable)
	r.POST("/api/v1/workspace/update", workspaceUpdate)
	r.POST("/api/v1/docker/reset", dockerReset)

	// proxy
	r.GET("/api/v1/config/proxy", gitProxyGet)
	r.POST("/api/v1/config/proxy", gitProxySet)
	r.GET("/api/v1/proxy/heliumApi", heliumApiProxy)

	r.GET("/state", stateHandler)
	r.GET("/registInfo", registInfoHandler)
}

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

func registInfoHandler(c *gin.Context) {
	d, err := register.GetRegistInfo()
	if err != nil {
		c.JSON(500, RespBody{Code: 500, Message: "error: " + err.Error()})
	} else {
		c.JSON(200, RespOK(d))
	}
}
