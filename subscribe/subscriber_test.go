package subscribe

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/ml444/scheduler/subscribe/call"
	"reflect"
	"testing"
)

func TestNewSubscriber(t *testing.T) {
	type args struct {
		namespace string
		topicName string
		subCfg    *Config
	}
	tests := []struct {
		name string
		args args
		want *Subscriber
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSubscriber(tt.args.namespace, tt.args.topicName, tt.args.subCfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSubscriber() = %v, want %v", got, tt.want)
			}
		})
	}
}

type rsp struct {
	msg  string
	code int32
}

func (r *rsp) Reset()         {}
func (r *rsp) String() string { return "" }
func (r *rsp) ProtoMessage()  {}
func TestSubscriber_NewResponse(t *testing.T) {

	type fields struct {
		Name          string
		Topic         string
		Route         string
		Addrs         []string
		Cfg           *Config
		Request       proto.Message
		Response      proto.Message
		BeforeProcess func(ctx context.Context, meta *call.MsgMeta)
		AfterProcess  func(ctx context.Context, meta *call.MsgMeta, req, rsp interface{}) (isRetry bool, ignoreRetryCount bool)
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			name:   "test",
			fields: fields{Response: &rsp{msg: "successfully", code: 1234}},
			want:   &rsp{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Subscriber{
				Response: tt.fields.Response,
			}
			if got := s.NewResponse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

type req struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Age  uint32 `protobuf:"varint,2,opt,name=age" json:"age,omitempty"`
}

func (r *req) Reset()         {}
func (r *req) String() string { return "" }
func (r *req) ProtoMessage()  {}
func TestSubscriber_UnMarshalRequest(t *testing.T) {

	r := &req{Name: "bar", Age: 456}
	bb, err := proto.Marshal(r)
	if err != nil {
		return
	}
	//t.Log(len(bb))
	type fields struct {
		Name          string
		Topic         string
		Route         string
		Addrs         []string
		Cfg           *Config
		Request       proto.Message
		Response      proto.Message
		BeforeProcess func(ctx context.Context, meta *call.MsgMeta)
		AfterProcess  func(ctx context.Context, meta *call.MsgMeta, req, rsp interface{}) (isRetry bool, ignoreRetryCount bool)
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				Request: &req{
					//Id: "foo",
					//Age:  123,
				},
			},
			//args: args{data: []byte(`{"name": "bar", "age": 456}`)},
			args: args{data: bb},
			want: &req{
				Name: "bar",
				Age:  456,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Subscriber{
				Id:            tt.fields.Name,
				Topic:         tt.fields.Topic,
				Route:         tt.fields.Route,
				Addrs:         tt.fields.Addrs,
				Cfg:           tt.fields.Cfg,
				Request:       tt.fields.Request,
				Response:      tt.fields.Response,
				BeforeProcess: tt.fields.BeforeProcess,
				AfterProcess:  tt.fields.AfterProcess,
			}
			got, err := s.UnMarshalRequest(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnMarshalRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnMarshalRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
