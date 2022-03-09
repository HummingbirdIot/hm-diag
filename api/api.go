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

func Route(r *gin.Engine, webFiles embed.FS, swagFiles embed.FS) {
	webFS = webFiles
	diagTask = diag.TaskInstance()
	register = regist.Instance()

	// web static files
	d, _ := fs.Sub(webFS, "web/release")
	r.StaticFS("/web", http.FS(d))
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/web")
	})

	d, _ = fs.Sub(swagFiles, "api/swagger_ui")
	r.StaticFS("/docs", http.FS(d))
	r.StaticFS("/docs-inner", http.FS(d))

	// hotspot
	// swagger:route GET /api/v1/device/state public device-state
	// Get device info
	//
	// this will show device state
	// Responses:
	//   200:DeviceInfo
	r.GET("/api/v1/device/state", deviceInfo)

	// swagger:route POST /api/v1/device/reboot public device-reboot
	// Reboot Device
	//
	// API will return immediately, you can check
	// Responses:
	//   200:EmptyBody
	r.POST("/api/v1/device/reboot", rebootDevice)

	// swagger:route POST /api/v1/device/light/blink inner device-light
	// Set Device light color
	//
	// Responses:
	//   200:EmptyBody
	r.POST("/inner/api/v1/device/light/blink", deviceLightBlink)

	// swagger:route GET /api/v1/lan/hotspot public lan-hotspots
	// Get devices(hotspots) address in LAN
	//
	// Device discovery by net interface `eth0`
	// Responses:
	//   200:DevDis
	r.GET("/api/v1/lan/hotspot", lanHotspotsLoad)

	// miner
	// swagger:route POST /api/v1/miner/resync public miner-resync
	// Resync miner
	//
	// Clean miner data and restart miner, miner will resync data
	// Responses:
	//   200:EmptyBody
	r.POST("/api/v1/miner/resync", resyncMiner)

	// swagger:route POST /api/v1/miner/restart public miner-restart
	// Restart miner
	//
	// Restart miner container
	// Responses:
	//   200:EmptyBody
	r.POST("/api/v1/miner/restart", restartMiner)

	// swagger:route GET /api/v1/miner/state public miner-state
	// Get miner info
	//
	// Get miner info and state
	// Responses:
	//   200:MinerInfo
	r.GET("/api/v1/miner/state", minerInfo)

	// swagger:route POST /inner/api/v1/miner/onboarding/txn inner miner-onboarding-txn
	// Generate onboarding transaction
	//
	// Invoke miner to generate onboarding transaction
	// Responses:
	r.POST("/inner/api/v1/miner/onboarding/txn", genOnboardingTxn)

	// swagger:route POST /inner/api/v1/miner/snapshot inner miner-snapshot
	// Take miner snapshot
	//
	// Call miner to take snapshot and return immediately.
	// You can check if snapshot taking is done by api `/inner/api/v1/miner/snapshot/state`
	// Responses:
	//   200:EmptyBody
	r.POST("/inner/api/v1/miner/snapshot", snapshotTake)

	// swagger:route GET /inner/api/v1/miner/snapshot/state inner miner-snapshot-state
	// Get state of snapshot taking
	//
	// Get the latest snapshot info
	// Responses:
	//   200:SnapshotStateRes
	r.GET("/inner/api/v1/miner/snapshot/state", snapshotState)

	// swagger:route GET /inner/api/v1/miner/snapshot/file/{name} inner miner-snapshot-download
	// Download snapshot file
	//
	// Snapshot file should be exist before call this api. you can call `/inner/api/v1/miner/snapshot/state` to see
	// Responses:
	//   200:EmptyBody
	r.GET("/inner/api/v1/miner/snapshot/file/:name", snapshotDownload)

	// swagger:route POST /inner/api/v1/miner/snapshot/apply inner miner-snapshot-apply
	// Apply snapshot
	//
	// Upload snapshot file and apply it
	// Responses:
	//   200:EmptyBody
	r.POST("/inner/api/v1/miner/snapshot/apply", snapshotApply)

	// swagger:route GET /inner/api/v1beta/miner/log inner miner-log
	// Query miner log
	//
	// Query miner log
	// Responses:
	//   200:StringBody
	r.GET("/inner/api/v1/log", logQuery)

	// workspace
	// swagger:route POST /inner/api/v1/workspace/reset inner workspace-reset
	// Reset workspace
	//
	// Pull new worksapce (main git repo) from server, delete old worksapce
	// Responses:
	//   200:EmptyBody
	r.POST("/inner/api/v1/workspace/reset", workspaceReset)

	// swagger:route GET /api/v1/workspace/update public workspace-update-get
	// Check workspace update
	//
	// Whether worksapce (main git repo) is update available
	// Responses:
	//   200:BoolBody
	r.GET("/api/v1/workspace/update", workspaceUpdateAvailable)

	// swagger:route POST /api/v1/workspace/update public workspace-update
	// Update worksapce (main git repo)
	//
	// Trigger workspace update and return immediately
	// Responses:
	//   200:EmptyBody
	r.POST("/api/v1/workspace/update", workspaceUpdate)

	// swagger:route POST /inner/api/v1/docker/reset inner docker-reset
	// Reset docker
	//
	// Delete all state  of docker (contain images、containers)
	// Responses:
	//   200:EmptyBody
	r.POST("/inner/api/v1/docker/reset", dockerReset)

	// proxy
	// swagger:route GET /api/v1/config/proxy public config-proxy-get
	// Get proxy config
	//
	// Proxy config is about git repository or git release files
	// `item` query parameter shoulbe: "gitRelease" or "gitRepo"
	//
	// responses:
	//	200:ProxyItem
	r.GET("/api/v1/config/proxy", gitProxyGet)

	// swagger:route POST /api/v1/config/proxy public config-proxy-update
	// Set proxy config
	//
	// roxy config is about git repository or git release files
	// `item` query parameter shoulbe: "gitRelease" or "gitRepo"
	// Responses:
	//   200:EmptyBody
	r.POST("/api/v1/config/proxy", gitProxySet)

	// swagger:route GET /inner/api/v1/proxy/heliumApi inner proxy-heliumApi
	// Proxy Helium API
	//
	// The Helium API uses the HTTPS protocol,
	// but some browsers do not allow access to the HTTPS API in HTTP sites,
	// so access through this API proxy it
	//
	// Responses:
	//   200:EmptyBody
	r.GET("/inner/api/v1/proxy/heliumApi", heliumApiProxy)

	// swagger:route GET /inner/state inner all-state
	// Get all state
	//
	// Get miner state and device state
	//
	// Responses:
	//   200:AllState
	r.GET("/inner/state", stateHandler)

	// TODO remove this route after next two version
	r.GET("/state", stateHandler)
	r.GET("/inner/registInfo", registInfoHandler)
	r.GET("/inner/api/v1/pktfwd/state", pktfwdVersion)
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
