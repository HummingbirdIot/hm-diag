package grpc

import (
	"context"
	"time"

	"github.com/kpango/glg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "xdt.com/hm-diag/diag/grpc/pb"
)

type Client struct {
	Url string
}

func (c Client) GatewayInfo() (res *pb.HeightRes, err error) {
	glg.Info("grpc URL: ", c.Url)
	conn, err := grpc.Dial(c.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glg.Errorf("did not connect: %v", err)
		return nil, err
	}
	defer conn.Close()
	co := pb.NewApiClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	r, err := co.Height(ctx, &pb.HeightReq{})
	if err != nil {
		glg.Errorf("could not greet: %v", err)
		return nil, err
	}
	glg.Infof("%+v", r)
	glg.Info(string(r.Gateway.Address))
	return r, nil
}

func (c Client) AddMinerGateway(owner []byte, payer []byte, stakingMode pb.GatewayStakingMode) (res string, err error) {
	glg.Info(owner, payer)
	glg.Info(string(owner), string(payer))
	conn, err := grpc.Dial(c.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glg.Errorf("did not connect: %v", err)
		return "", err
	}
	defer conn.Close()
	co := pb.NewApiClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	r, err := co.AddGateway(ctx, &pb.AddGatewayReq{Owner: owner, Payer: payer, StakingMode: stakingMode})
	if err != nil {
		glg.Errorf("could not greet: %v", err)
		return "", err
	}
	glg.Infof("%+v", r)
	glg.Info(string(r.AddGatewayTxn))
	return string(r.AddGatewayTxn), nil
}
