package api

import (
	"encoding/base64"
	"log"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"xdt.com/hm-diag/ctrl"
	"xdt.com/hm-diag/diag"
	"xdt.com/hm-diag/diag/miner"
)

func resyncMiner(c *gin.Context) {
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
}

func restartMiner(c *gin.Context) {
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
}

func snapshotTake(c *gin.Context) {
	log.Println("to snapshot")
	ctrl.SnapshotTake()
	c.JSON(200, RespOK(nil))
}

func snapshotState(c *gin.Context) {
	log.Println("to get snapshot state")
	res, err := ctrl.SnapshotState()
	if err != nil {
		log.Printf("get snapshot state error: %+v\n", err)
		c.JSON(500, RespBody{Code: 500, Message: err.Error()})
	} else {
		c.JSON(200, RespOK(res))
	}
}

func snapshotDownload(c *gin.Context) {
	log.Println("to get snapshot file")
	f := c.Param("name")
	b, err := base64.StdEncoding.DecodeString(f)
	if err != nil {
		log.Println("get snapshot, wrong path param", f)
		c.JSON(400, RespBody{Code: 400, Message: "wrong path param"})
		return
	}
	f = string(b)
	if !strings.HasPrefix(f, "/tmp/") {
		log.Println("get snapshot, wrong path param", f)
		c.JSON(400, RespBody{Code: 400, Message: "wrong path param"})
		return
	}
	c.Header("Content-Disposition", "attachment; filename=\""+path.Base(f)+"\"")
	c.File(f)
}

func snapshotApply(c *gin.Context) {
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
}

func minerInfo(c *gin.Context) {
	var d diag.TaskData
	c.Query("cache")
	if c := c.Query("cache"); c == "true" {
		d = diagTask.MinerInfo()
	} else {
		d = diag.TaskData{Data: diagTask.FetchMinerInfo(), FetchTime: time.Now()}
	}
	c.JSON(200, RespOK(d.Data))
}

func minerLogQuery(c *gin.Context) {
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