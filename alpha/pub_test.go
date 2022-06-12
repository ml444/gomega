package alpha

import (
	"context"
	"testing"

	"github.com/ml444/gomega/alpha/omega"
)

func TestPub(t *testing.T) {
	rsp, err := Pub(context.Background(), &omega.PubReq{
		Topic:    "test_topic",
		HashCode: 123456,
		Data:     []byte("test1"),
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(rsp)
}
