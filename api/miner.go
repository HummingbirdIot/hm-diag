package api

import (
	"encoding/base64"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
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

func genOnboardingTxn(c *gin.Context) {
	ownerAddr := c.Query("owner")
	glg.Info("to generate onboarding txn for owner:", ownerAddr)
	if ownerAddr == "" {
		c.JSON(400, RespBody{Code: 400, Message: "owner address must be provided"})
		return
	}
	txn, err := ctrl.GenOnboardingTxn(ownerAddr, ctrl.MakerAddr)
	if err != nil {
		glg.Error("generating onboarding txn error", err)
		c.JSON(500, RespBody{Code: 500, Message: err.Error()})
		return
	}
	c.JSON(200, RespOK(txn))
}

func genAssertLocationTxn(c *gin.Context) {
	ownerAddr := c.Query("owner")
	payerAddr := c.Query("payer")
	location := c.Query("h3")
	nonceStr := c.Query("nonce")
	nonce, err := strconv.Atoi(nonceStr)
	if err != nil {
		nonce = 1
	}
	glg.Infof("to generate assert location txn for owner: %s, payer: %s", ownerAddr, payerAddr)
	if ownerAddr == "" || payerAddr == "" || location == "" {
		c.JSON(400, RespBody{Code: 400, Message: "owner, payer and h3 must be provided"})
		return
	}
	txn, err := ctrl.GenAssertLocationTxn(ownerAddr, payerAddr, location, nonce)
	if err != nil {
		glg.Info("generating assert location txn error:", err)
		c.JSON(500, RespBody{Code: 500, Message: err.Error()})
		return
	}
	c.JSON(200, RespOK(txn))
}

func snapshotTake(c *gin.Context) {
	glg.Debug("to snapshot")
	ctrl.SnapshotTake()
	c.JSON(200, RespOK(nil))
}

func snapshotState(c *gin.Context) {
	glg.Debug("to get snapshot state")
	res, err := ctrl.SnapshotState()
	if err != nil {
		glg.Errorf("get snapshot state error: %+v\n", err)
		c.JSON(500, RespBody{Code: 500, Message: err.Error()})
	} else {
		c.JSON(200, RespOK(res))
	}
}

func snapshotDownload(c *gin.Context) {
	glg.Debug("to get snapshot file")
	f := c.Param("name")
	b, err := base64.StdEncoding.DecodeString(f)
	if err != nil {
		glg.Errorf("get snapshot, wrong path param", f)
		c.JSON(400, RespBody{Code: 400, Message: "wrong path param"})
		return
	}
	f = string(b)
	if !strings.HasPrefix(f, "/tmp/") {
		glg.Errorf("get snapshot, wrong path param", f)
		c.JSON(400, RespBody{Code: 400, Message: "wrong path param"})
		return
	}
	c.Header("Content-Disposition", "attachment; filename=\""+path.Base(f)+"\"")
	c.File(f)
}

func snapshotApply(c *gin.Context) {
	glg.Debug("to apply snapshot")
	f, err := c.FormFile("file")
	if err != nil {
		glg.Errorf("get form file error", err)
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

func pktfwdVersion(c *gin.Context) {
	v, err := miner.PacketForwardVersion()
	if err == nil {
		c.JSON(200, RespOK(map[string]interface{}{
			"version": v,
		}))
	} else {
		c.JSON(500, RespBody{Code: 500, Message: err.Error()})
	}
}
