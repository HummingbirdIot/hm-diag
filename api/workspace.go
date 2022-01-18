package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"xdt.com/hm-diag/ctrl"
)

func dockerReset(c *gin.Context) {
	ctrl.DockerReset()
	c.JSON(200, RespOK(nil))
}

func workspaceReset(c *gin.Context) {
	err := ctrl.GitRepoReset()
	if err != nil {
		log.Println(err)
		c.JSON(500, RespBody{Code: 500, Message: err.Error()})
	} else {
		c.JSON(200, RespOK(nil))
	}
}

func workspaceUpdateAvailable(c *gin.Context) {
	r, err := ctrl.IsGitRepoToUpdate()
	if err != nil {
		c.JSON(500, RespBody{Code: 500, Message: err.Error()})
	} else {
		c.JSON(200, RespOK(r))
	}
}

func workspaceUpdate(c *gin.Context) {
	ctrl.ExecMainUpdate()
	c.JSON(200, RespOK(nil))
}
