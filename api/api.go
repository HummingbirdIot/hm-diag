package api

import (
	"context"
	"embed"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"xdt.com/hm-diag/config"
	"xdt.com/hm-diag/diag"
	"xdt.com/hm-diag/diag/device"
	"xdt.com/hm-diag/diag/miner"
	"xdt.com/hm-diag/link"
	"xdt.com/hm-diag/util"
)

type RespBody struct {
	Data    interface{} `json:"data"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
}

type PasswordBody struct {
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}

type ClientConfigBody struct {
	Secret string `json:"secret"`
	Server string `json:"server"`
}

func RespOK(data interface{}) RespBody {
	return RespBody{Data: data, Code: 200, Message: "OK"}
}

var diagTask *diag.Task
var webFS embed.FS

func Route(r *gin.Engine, webFiles embed.FS, swagFiles embed.FS) {
	webFS = webFiles
	diagTask = diag.TaskInstance()

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
	r.POST("/inner/api/v1/miner/txn/onboarding", genOnboardingTxn)

	r.POST("/inner/api/v1/miner/txn/assertLocation", genAssertLocationTxn)

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

	r.POST("/inner/api/v1/config/safe", saveConfigHandle)
	r.GET("/inner/api/v1/config/safe", getConfigHandle)
	r.GET("/inner/api/v1/safe/isViaPrivate", isViaPrivate)

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
	r.GET("/inner/state", stateInnerHandler)

	r.GET("/inner/api/v1/version", versionHandler)

	// TODO remove this route after next two version
	r.GET("/state", stateHandler)
	r.GET("/inner/api/v1/pktfwd/state", pktfwdVersion)
	r.POST("/inner/api/v1/clientConfig/safe", saveClientConfigHandle)

	r.GET("/inner/api/v1/network/ping", networkTestHandler)
	r.POST("/api/v1/login", loginHandler)
	r.POST("/api/v1/password", passwordHandler)
}

func stateInnerHandler(c *gin.Context) {
	v, err := miner.PacketForwardVersion()
	if err != nil {
		c.JSON(500, RespBody{Code: 500, Message: err.Error()})
		return
	}
	var res map[string]interface{}
	if c := c.Query("cache"); c == "true" {
		res = diagTask.Data().Data
		res["time"] = diagTask.Data().FetchTime
		res["packetForwardVersion"] = v
	} else {
		res = map[string]interface{}{
			"device":               diagTask.FetchDeviceInfo(),
			"miner":                diagTask.FetchMinerInfo(),
			"packetForwardVersion": v,
		}
	}
	res["notice"] = `do not use this api path "/" to integrate, use api under path "api/"`
	c.JSON(200, RespOK(res))
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

func saveClientConfigHandle(c *gin.Context) {
	var conf ClientConfigBody
	err := json.NewDecoder(c.Request.Body).Decode(&conf)
	if err != nil {
		log.Println(err)
		c.JSON(400, RespBody{Code: 400, Message: "invalid request body for config"})
		return
	}
	existConfig, err := link.LoadClientConfig()
	if err != nil {
		log.Println("get clientConfig error : ", err)
		c.JSON(500, RespBody{Code: 500, Message: "get clientConfig error : " + err.Error()})
		return
	}
	var newConfig link.ClientConfig
	newConfig.ID = existConfig.ID
	newConfig.Secret = conf.Secret
	newConfig.Server = conf.Server
	log.Printf("to save clent config file, content: %#v", conf)
	err = link.SaveClientConfig(newConfig)
	if err != nil {
		log.Println("save clent config file error:", err)
		c.JSON(500, RespBody{Code: 500, Message: err.Error()})
		return
	}
	log.Printf("saved client config file, content: %#v", newConfig)
	log.Printf("to reconnect link server: %#v", conf.Server)
	err = link.Start(context.Background())
	if err != nil {
		log.Println("reconnect link error : ", err)
		c.JSON(500, RespBody{Code: 500, Message: "save client config successful, but reconnect error : " + err.Error()})
		return
	}
	log.Printf("reconnected link server: %#v", conf.Server)
	c.JSON(200, RespOK(nil))
}

func saveConfigHandle(c *gin.Context) {
	var conf config.ConfiFileData
	err := json.NewDecoder(c.Request.Body).Decode(&conf)
	if err != nil {
		c.JSON(400, RespBody{Code: 400, Message: "invalid request body for config"})
		return
	}
	log.Printf("to save config file, content: %#v", conf)
	err = config.SaveConfigFile(conf)
	if err != nil {
		log.Println("save config file error:", err)
		c.JSON(500, RespBody{Code: 500, Message: err.Error()})
		return
	}
	log.Printf("saved config file, content: %#v", conf)
	c.JSON(200, RespOK(nil))
}

func getConfigHandle(c *gin.Context) {
	conf, err := config.ReadConfigFile()
	if err != nil {
		log.Println("get config file error:", err)
		c.JSON(500, RespBody{Code: 500, Message: err.Error()})
		return
	}

	var res = make(map[string]interface{})
	res["publicAccess"] = conf.PublicAccess
	c.JSON(200, RespOK(res))
}

func versionHandler(c *gin.Context) {
	c.JSON(200,
		RespOK(gin.H{
			"version": config.Version,
			"githash": config.Githash,
		}))
}

func isViaPrivate(c *gin.Context) {
	rIp, ok := c.RemoteIP()
	if !ok {
		c.JSON(400, RespBody{Code: 400, Message: "can't get remote IP"})
		return
	}
	r := util.IsPrivateIp(rIp)
	c.JSON(200, RespOK(r))
}

func networkTestHandler(c *gin.Context) {
	d := c.Query("ip")
	err := device.NetworkTest(d)
	if err != nil {
		c.JSON(500, RespBody{Code: 500, Message: err.Error()})
		return
	}
	c.JSON(200, RespOK(nil))
}

func loginHandler(c *gin.Context) {
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Writer.WriteHeader(500)
		return
	}
	if config.Config().Password == string(b) {
		c.JSON(200, RespOK(GenToken()))
	} else {
		c.Writer.WriteHeader(700)
	}
}

func passwordHandler(c *gin.Context) {
	var requestData PasswordBody
	err := json.NewDecoder(c.Request.Body).Decode(&requestData)
	if err != nil {
		c.JSON(400, RespBody{Code: 400, Message: "invalid request body for password"})
		return
	}
	confFile := config.Config()

	if confFile.Password == requestData.Password {
		var newConfig config.ConfiFileData
		newConfig.Password = requestData.NewPassword
		newConfig.PublicAccess = confFile.PublicAccess
		err = config.SaveConfigFile(newConfig)
		if err != nil {
			log.Println("method SetPassword error:", err)
			c.JSON(500, RespBody{Code: 500, Message: err.Error()})
			return
		}
		c.JSON(200, RespOK("ok"))
	} else {
		c.Writer.WriteHeader(700)
	}
}
