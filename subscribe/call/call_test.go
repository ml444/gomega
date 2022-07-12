package call

import (
	"context"
	"google.golang.org/protobuf/runtime/protoimpl"
	"testing"
)

type Point struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Latitude  int32 `protobuf:"varint,1,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude int32 `protobuf:"varint,2,opt,name=longitude,proto3" json:"longitude,omitempty"`
}

func (x *Point) Reset() {
	*x = Point{}
}

func (x *Point) String() string {
	return ""
}

func (*Point) ProtoMessage() {}


func (x *Point) GetLatitude() int32 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *Point) GetLongitude() int32 {
	if x != nil {
		return x.Longitude
	}
	return 0
}
type Feature struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the feature.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The point where the feature is detected.
	Location *Point `protobuf:"bytes,2,opt,name=location,proto3" json:"location,omitempty"`
}

func (x *Feature) Reset() {
	*x = Feature{}
}

func (x *Feature) String() string {
	return ""
}

func (*Feature) ProtoMessage() {}

func TestCall(t *testing.T) {
	type args struct {
		ctx     context.Context
		route   string
		in      interface{}
		out     interface{}
		timeout uint32
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test_GetFeature",
			args:    args{
				ctx:     context.TODO(),
				route:   "routeguide.RouteGuide/GetFeature",
				in:      &Point{Latitude: 409146138, Longitude: -746188906},
				out:     &Feature{},
				timeout: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Call(tt.args.ctx, tt.args.route, tt.args.in, tt.args.out, tt.args.timeout); (err != nil) != tt.wantErr {
				t.Errorf("Call() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
