package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"xdt.com/hm-diag/ctrl"
	"xdt.com/hm-diag/devdis"
	"xdt.com/hm-diag/diag"
)

func rebootDevice(c *gin.Context) {
	c.JSON(200, RespBody{
		Code:    200,
		Message: "receive reboot request, to reboot",
	})
	go ctrl.RebootDevice()
}

func deviceInfo(c *gin.Context) {
	var d diag.TaskData
	if c := c.Query("cache"); c == "true" {
		d = diagTask.DeviceInfo()
	} else {
		d = diag.TaskData{Data: diagTask.FetchDeviceInfo(), FetchTime: time.Now()}
	}
	c.JSON(200, RespOK(d.Data))
}

func lanHotspotsLoad(c *gin.Context) {
	c.JSON(http.StatusOK, RespOK(devdis.Services()))
}
