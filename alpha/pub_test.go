package alpha

import (
	"context"
	"testing"
)

func TestPub(t *testing.T) {
	rsp, err := Pub(context.Background(), &PubReq{
		Namespace:            "default",
		TopicName:            "test_topic",
		Partition:            0,
		HashCode:             123456,
		Data:                 []byte{},
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(rsp)
}
