package alpha

import (
	"context"
	"io"
	"time"
)

func Sub(ctx context.Context, req *SubReq) (*SubRsp, error) {
	var err error
	var rsp SubRsp
	Init()
	//return cli.Pub(ctx, req)
	err = conn.Invoke(ctx, "/omega.OmegaService/Sub", req, &rsp)
	if err != nil {
		return nil, err
	}
	return &rsp, nil
}
var client OmegaServiceClient
func init() {
	Init()
	client = NewOmegaServiceClient(conn)
}
func LoopConsume(handler func (r *ConsumeRsp) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	stream, err := client.Consume(ctx)
	if err != nil {
		return err
	}
	var rsp *ConsumeRsp
	req := ConsumeReq{
		Token:            "test_token",
		Sequence:         1,
		IsRetry:          false,
		IgnoreRetryCount: false,
		RetryIntervalMs:  0,
	}
	for {
		if err = stream.Send(&req); err != nil {
			return err
		}
		rsp, err = stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			println(err)
			return err
		}
		err = handler(rsp)
		if err != nil {
			println(err)
			return err
		}
		req.Sequence++
	}
}
