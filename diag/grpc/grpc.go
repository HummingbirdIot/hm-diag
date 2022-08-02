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

func (c Client) Height() (res uint64, err error) {
	glg.Info("grpc URL: ", c.Url)
	conn, err := grpc.Dial(c.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glg.Errorf("did not connect: %v", err)
		return 0, err
	}
	defer conn.Close()
	co := pb.NewApiClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	r, err := co.Height(ctx, &pb.HeightReq{})
	if err != nil {
		glg.Errorf("could not greet: %v", err)
		return 0, err
	}
	return r.Height, nil
}
