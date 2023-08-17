package alpha

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ml444/gomega/alpha/omega"
)

//type Point struct {
//	state         protoimpl.MessageState
//	sizeCache     protoimpl.SizeCache
//	unknownFields protoimpl.UnknownFields
//
//	Latitude  int32 `protobuf:"varint,1,opt,name=latitude,proto3" json:"latitude,omitempty"`
//	Longitude int32 `protobuf:"varint,2,opt,name=longitude,proto3" json:"longitude,omitempty"`
//}
//
//func (x *Point) Reset() {
//	*x = Point{}
//}
//
//func (x *Point) String() string {
//	return ""
//}
//
//func (*Point) ProtoMessage() {}
//
//
//func (x *Point) GetLatitude() int32 {
//	if x != nil {
//		return x.Latitude
//	}
//	return 0
//}
//
//func (x *Point) GetLongitude() int32 {
//	if x != nil {
//		return x.Longitude
//	}
//	return 0
//}
//type Feature struct {
//	state         protoimpl.MessageState
//	sizeCache     protoimpl.SizeCache
//	unknownFields protoimpl.UnknownFields
//
//	// The name of the feature.
//	Id string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
//	// The point where the feature is detected.
//	Location *Point `protobuf:"bytes,2,opt,name=location,proto3" json:"location,omitempty"`
//}
//
//func (x *Feature) Reset() {
//	*x = Feature{}
//}
//
//func (x *Feature) String() string {
//	return ""
//}
//
//func (*Feature) ProtoMessage() {}

//	func SendCall() {
//		ctx := context.TODO()
//		route := "routeguide.RouteGuide/GetFeature"
//		in := &Point{Latitude: 409146138, Longitude: -746188906}
//		out := &Feature{}
//		err := call.Call(ctx, route, in, out, 123)
//		if err != nil {
//			println(err)
//		}
//	}
func TestSubscribe(t *testing.T) {
	type args struct {
		ctx context.Context
		req *omega.SubReq
	}
	tests := []struct {
		name    string
		args    args
		want    *omega.SubRsp
		wantErr bool
	}{
		{
			name: "",
			args: args{
				ctx: context.TODO(),
				req: &omega.SubReq{
					ClientId: "test_client_id",
					Topic:    "test_topic",
					Group:    "test_group",
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sub(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Subscribe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Subscribe() got = %v, want %v", got, tt.want)
			}
			t.Log(got)
		})
	}
}

var f *os.File

func handler(ctx context.Context, r *omega.ConsumeRsp) error {
	f.Write(r.Data)
	f.Write([]byte("\n"))
	fmt.Println(string(r.Data), r.Sequence)
	return nil
}

func TestLoopConsume(t *testing.T) {
	InitClient("localhost:12345")
	var err error
	f, err = os.OpenFile("test1.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	type args struct {
		ctx      context.Context
		topic    string
		group    string
		clientId string
		handler  consumeFunc
	}
	ctx := context.TODO()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"", args{ctx: ctx, topic: "test_topic", group: "test_group", clientId: "test_client_id", handler: handler}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			if err := LoopConsume(tt.args.ctx, tt.args.topic, tt.args.group, tt.args.clientId, tt.args.handler); (err != nil) != tt.wantErr {
				t.Errorf("LoopConsume() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Log("===>time:", time.Now().Sub(start).Seconds())
		})
	}
}

//func TestSub(t *testing.T) {
//	subRsp, err := Sub(context.TODO(), &omega.SubReq{
//		ClientId: "test_client_id",
//		Topic:    "test_topic",
//		Group:    "test_group",
//	})
//	if err != nil {
//		t.Errorf("err: %v", err)
//		return
//	}
//}
