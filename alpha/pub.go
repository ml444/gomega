package alpha

import (
	"context"

	"github.com/ml444/gomega/alpha/omega"
)

func Pub(ctx context.Context, req *omega.PubReq) (*omega.Response, error) {
	var err error
	if cli == nil {
		err = InitClient("localhost:12345")
		if err != nil {
			return nil, err
		}
	}
	return cli.Pub(ctx, req)
}
