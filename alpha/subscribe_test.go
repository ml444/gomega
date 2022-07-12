package alpha

import (
	"context"
	"github.com/ml444/scheduler/subscribe/call"
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
//	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
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

func SendCall() {
	ctx := context.TODO()
	route := "routeguide.RouteGuide/GetFeature"
	in := &Point{Latitude: 409146138, Longitude: -746188906}
	out := &Feature{}
	err := call.Call(ctx, route, in, out, 123)
	if err != nil {
		println(err)
	}
}
func TestSubscribe(t *testing.T) {
	type args struct {
		ctx context.Context
		req *SubscribeReq
	}
	tests := []struct {
		name    string
		args    args
		want    *SubscribeRsp
		wantErr bool
	}{
		{
			name:    "",
			args:    args{
				ctx: context.TODO(),
				req: &SubscribeReq{
					Namespace:            "",
					Name:                 "",
					Topic:                "",
					Route:                "routeguide.RouteGuide/GetFeature",
					Request:              &Point{},
					Response:             &Feature{},
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Subscribe(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Subscribe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Subscribe() got = %v, want %v", got, tt.want)
			}
		})
	}
}
