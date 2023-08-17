package alpha

import (
	"context"
	"testing"

	"github.com/ml444/gomega/alpha/omega"
)

func TestPub(t *testing.T) {
	for i := 0; i < 10000; i++ {
		_, err := Pub(context.Background(), &omega.PubReq{
			Topic:    "test_topic",
			HashCode: 123456,
			Data:     []byte(`{"test": "balabala", "far": "away"}`),
		})
		if err != nil {
			t.Error(err)
		}
		//time.Sleep(time.Millisecond * 10)
		t.Log(i)
	}
	//rsp, err := Pub(context.Background(), &omega.PubReq{
	//	Topic:    "test_topic",
	//	HashCode: 123456,
	//	Data:     []byte("test1"),
	//})
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Log(rsp)
}
