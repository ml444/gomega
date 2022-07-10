package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/ml444/scheduler/brokers"
	"github.com/ml444/scheduler/config"
	"github.com/ml444/scheduler/publish"
	"google.golang.org/grpc"
	"net"
	"time"
)

type OmegaServer struct {
}

func (s *OmegaServer) Pub(ctx context.Context, req *PubRequest) (*PubResponse, error) {
	var rsp PubResponse
	fmt.Println(req)

	dataSize := len(req.Data)
	if dataSize > config.MaxDataSize {
		return nil, errors.New("data length exceeds limit")
	}
	item := brokers.Item{
		CreatedAt:  uint32(time.Now().Unix()),
		HashCode:   req.HashCode,
		Partition:  req.Partition,
		Size:       uint32(dataSize),
		//DelayType:  req.DelayType,
		//DelayValue: req.DelayValue,
		//Priority:   req.Priority,
		Data:       req.Data,
		Namespace:  req.Namespace,
	}

	err := publish.Pub(&ctx, req.Namespace, req.TopicName, &item)
	if err != nil {
		return nil, err
	}
	return &rsp, nil
}

func RunServer() error {
	listen, err := net.Listen("tcp", ":12345")
	if err != nil {
		return err
	}
	svr := grpc.NewServer()
	RegisterOmegaServiceServer(svr, &OmegaServer{})
	err = svr.Serve(listen)
	if err != nil {
		return err
	}
	return nil
}
