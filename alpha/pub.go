package alpha

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var conn *grpc.ClientConn

//var cli OmegaServiceClient

func Init() {
	var err error
	serviceAddr := "127.0.0.1:12345"
	conn, err = grpc.Dial(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		println(err)
		return
	}
	//defer conn.Close()
	//cli = NewOmegaServiceClient(conn)

}
func Pub(ctx context.Context, req *PubRequest) (*PubResponse, error) {
	var err error
	var rsp PubResponse
	Init()
	//return cli.Pub(ctx, req)
	err = conn.Invoke(ctx, "/omega.OmegaService/Pub", req, &rsp)
	if err != nil {
		return nil, err
	}
	return &rsp, nil
}
