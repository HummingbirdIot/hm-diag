package link

import (
	"context"
	"log"
)

var (
	singleClient     *Client
	preLinkCtx       context.Context
	preLinkCtxCancel context.CancelFunc
)

func init() {
	singleClient = &Client{}
}

func Start(rootCtx context.Context) error {
	if preLinkCtxCancel != nil {
		preLinkCtxCancel()
	}
	ctx, ctxCancel := context.WithCancel(rootCtx)
	preLinkCtx, preLinkCtxCancel = ctx, ctxCancel

	conf, err := GetClientConfig()
	if err != nil {
		linkLog(err)
		return err
	}
	singleClient.config = *conf

	err = singleClient.Start(ctx)
	if err != nil {
		linkLog(err)
		return err
	}
	log.Println(">>>>>>>>> Link client connect success")

	return nil
}

func ReportData(data any) error {
	singleClient.WriteMessage(data)
	return nil
}

func linkLog(err error) {
	log.Println("link start failed", err)
}
