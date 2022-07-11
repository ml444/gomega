package call

import (
	"context"
	log "github.com/ml444/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MsgMeta struct {
	CreatedAt uint32 `json:"created_at,omitempty"`
	RetryCnt  uint32 `json:"retry_cnt,omitempty"`
	Data      []byte `json:"data,omitempty"`
	MsgId     string `json:"msg_id,omitempty"`
	// @desc: 服务端最大重试次数，业务可以根据这个来做最后失败逻辑
	//MaxRetryCnt uint32 `protobuf:"varint,7,opt,name=max_retry_cnt,json=maxRetryCnt" json:"max_retry_cnt,omitempty"`
	//Timeout     uint32
}

type Request struct {
	Route string
	In    interface{}
	Out   interface{}

	Timeout     uint32
}

func (r *Request) Do(ctx context.Context) error {
	conn, err := r.getConn()
	if err != nil {
		log.Error()
		return err
	}
	return conn.Invoke(ctx, r.Route, r.In, r.Out)
}

func (r *Request) getConn() (*grpc.ClientConn, error) {
	serviceAddr, err := r.getAddr("")
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return grpc.Dial(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func (r *Request) getAddr(key string) (string, error) {
	return "127.0.0.1:50051", nil
}
