package alpha

import (
	"context"
	"fmt"
	"reflect"
	"testing"
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

//func SendCall() {
//	ctx := context.TODO()
//	route := "routeguide.RouteGuide/GetFeature"
//	in := &Point{Latitude: 409146138, Longitude: -746188906}
//	out := &Feature{}
//	err := call.Call(ctx, route, in, out, 123)
//	if err != nil {
//		println(err)
//	}
//}
func TestSubscribe(t *testing.T) {
	type args struct {
		ctx context.Context
		req *SubReq
	}
	tests := []struct {
		name    string
		args    args
		want    *SubRsp
		wantErr bool
	}{
		{
			name: "",
			args: args{
				ctx: context.TODO(),
				req: &SubReq{
					ClientId:      "test_client_id",
					Namespace:     "default",
					Topic:         "test_topic",
					Group:         "test_group",
					Policy:        0,
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


func handler(r *ConsumeRsp) error {
	fmt.Println(string(r.Data))
	return nil
}

func TestLoopConsume(t *testing.T) {
	type args struct {
		handler func(r *ConsumeRsp) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test",
			args:    args{handler: handler},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoopConsume(tt.args.handler); (err != nil) != tt.wantErr {
				t.Errorf("LoopConsume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}