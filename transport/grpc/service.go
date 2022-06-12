package grpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/ml444/gomega/alpha/omega"
	"github.com/ml444/gomega/broker"
	"github.com/ml444/gomega/config"
	"github.com/ml444/gomega/log"
)

type OmegaServer struct {
	omega.UnsafeOmegaServiceServer
}

func generateRandomHashCode() uint64 {
	return rand.Uint64()
}

func (s *OmegaServer) Pub(ctx context.Context, req *omega.PubReq) (*omega.Response, error) {
	var rsp omega.Response

	dataSize := len(req.Data)
	if dataSize > config.GlobalCfg.MaxMessageSize {
		return nil, errors.New("data length exceeds limit")
	}
	hashCode := req.HashCode
	if hashCode == 0 {
		hashCode = generateRandomHashCode()
	}
	item := broker.Item{
		CreatedAt: uint32(time.Now().Unix()),
		HashCode:  hashCode,
		Size:      uint32(dataSize),
		//DelayType:  req.DelayType,
		//DelayValue: req.DelayValue,
		//Priority:   req.Priority,
		Data: req.Data,
		//Namespace: req.Namespace,
	}
	bk := broker.GetOrCreateBroker(req.Topic)
	err := bk.Send(&item)
	if err != nil {
		return nil, err
	}

	rsp.Status = 10000
	return &rsp, nil
}

func (s *OmegaServer) Sub(ctx context.Context, req *omega.SubReq) (*omega.SubRsp, error) {
	var rsp omega.SubRsp
	// init worker or get worker
	if req.ClientId == "" {
		req.ClientId = fmt.Sprintf("%s_%s_%d", req.Topic, req.Group, time.Now().UnixNano())
	}
	log.Infof("Subscribe: Topic: %s, Group: %s, ClientId: %s", req.Topic, req.Group, req.ClientId)
	b := broker.GetOrCreateBroker(req.Topic)
	worker := b.GetConsumeWorker(req.Topic, req.Group, req.ClientId)
	broker.SetWorker(req.ClientId, worker)
	rsp.Status = 1000
	rsp.Token = req.ClientId
	return &rsp, nil
}

func (s *OmegaServer) Consume(stream omega.OmegaService_ConsumeServer) error {
	req, err := stream.Recv()
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}
	if req.Token == "" {
		return errors.New("token must be not empty")
	}
	// get worker
	worker := broker.GetWorker(req.Token)
	log.Infof("Consume: Token: %s", req.Token)
	var exit bool
	for !exit {
		select {
		case <-stream.Context().Done():
			return nil
		default:
		}
		item := worker.NextItem(&exit)
		if item == nil {
			return errors.New("item is nil")
		}
		log.Debugf("Consume: Token: %s, Item: %v", req.Token, item)
		err = stream.Send(&omega.ConsumeRsp{
			Sequence: item.Sequence,
			Data:     item.Data,
		})
		if err != nil {
			return err
		}

	}
	return nil
}

func (s *OmegaServer) UpdateSubCfg(ctx context.Context, req *omega.SubCfg) (*omega.Response, error) {
	var rsp omega.Response
	//topicIns, err := topic.GetTopic(req.Namespace, req.Topic)
	//if err != nil {
	//	return nil, err
	//}
	//cfg := subscribe.SubConfig{
	//	//Type:                req.Policy,
	//	Namespace:           req.Namespace,
	//	Topic:               req.Topic,
	//	GroupName:           req.Group,
	//	MaxRetryCount:       req.MaxRetryCount,
	//	MaxTimeout:          req.MaxTimeout,
	//	RetryIntervalMs:     req.RetryIntervalMs,
	//	ItemLifetimeInQueue: req.ItemLifetimeInQueue,
	//}
	//topicIns.AddSubCfg(&cfg)
	return &rsp, nil
}
func RunServer() error {
	listen, err := net.Listen("tcp", ":12345")
	if err != nil {
		return err
	}
	svr := grpc.NewServer()
	omega.RegisterOmegaServiceServer(svr, &OmegaServer{})
	err = svr.Serve(listen)
	if err != nil {
		return err
	}
	return nil
}
