package call

import (
	"context"
)

func Call(ctx context.Context, route string, in, out interface{}, timeout uint32) error {
	req := &Request{
		In:      in,
		Out:     out,
		Timeout: timeout,
		Route:   route,
	}
	err := req.Do(ctx)
	if err != nil {
		return err
	}
	return nil
}
