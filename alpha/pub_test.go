package alpha

import (
	"context"
	"testing"
)

func TestPub(t *testing.T) {
	rsp, err := Pub(context.Background(), &PubRequest{
		Namespace:            "test_topic",
		TopicName:            "topic_test_1",
		Partition:            0,
		HashCode:             123456,
		Data:                 []byte{},
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(rsp)
}
