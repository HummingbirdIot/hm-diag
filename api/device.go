package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"xdt.com/hm-diag/ctrl"
	"xdt.com/hm-diag/devdis"
	"xdt.com/hm-diag/diag"
	"xdt.com/hm-diag/diag/device"
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

func logQuery(c *gin.Context) {
	logType := c.Query("type")
	if logType != string(device.MINER_LOG) && logType != string(device.PWT_FWD_LOG) {
		c.JSON(400, RespBody{Code: 400, Message: "type should be miner or pktfwd"})
		log.Printf("query log invalid type %s\n", logType)
		return
	}
	filter := c.Query("filter")
	var since time.Time = time.Now().Add(time.Minute * time.Duration(5))
	var until time.Time = time.Now()
	var maxLines uint = 2000
	if logType == string(device.PWT_FWD_LOG) {
		if st, err := time.Parse("2006-01-02T15:04:05", c.Query("since")); err == nil {
			since = st
		} else {
			log.Println("query log, convert since time error: ", err)
		}
		if tt, err := time.Parse("2006-01-02T15:04:05", c.Query("until")); err == nil {
			until = tt
		} else {
			log.Println("query log, convert until time error: ", err)
		}
		log.Printf("query log, type: %s,since: %s, until: %s, filter: %s\n",
			logType, since, until, filter)
	} else {
		if l := c.Query("limit"); l != "" {
			limit, err := strconv.ParseInt(l, 10, 16)
			if err != nil {
				c.JSON(400, RespBody{Code: 400, Message: err.Error()})
				return
			}
			maxLines = uint(limit)
		}
		log.Printf("query log, type: %s,filter: %s, limit lines: %d\n",
			logType, filter, maxLines)
	}

	var l string
	var err error
	if logType == string(device.PWT_FWD_LOG) {
		l, err = device.QueryPktfwdLog(since, until, filter)
	} else {
		l, err = device.QueryMinerLog(filter, maxLines)
	}
	if err == nil {
		c.JSON(200, RespOK(l))
	} else {
		c.JSON(500, RespBody{Code: 500, Message: err.Error()})
		log.Println("query log error: ", err)
	}
}
