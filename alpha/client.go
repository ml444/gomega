package alpha

import (
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ml444/gomega/alpha/omega"
)

var conn *grpc.ClientConn
var cli omega.OmegaServiceClient

func InitClient(addr string) error {
	var err error
	if addr == "" {
		addr = os.Getenv("OMEGA_ADDR")
	}
	conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		println(err)
		return err
	}
	cli = omega.NewOmegaServiceClient(conn)
	return nil
}

func CloseClient() error {
	return conn.Close()
}
