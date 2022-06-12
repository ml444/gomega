package grpc

import (
	"context"
	"reflect"
	"testing"

	"github.com/ml444/gomega/alpha/omega"
)

func TestOmegaServer_Consume(t *testing.T) {
	type fields struct {
		UnsafeOmegaServiceServer omega.UnsafeOmegaServiceServer
	}
	type args struct {
		stream omega.OmegaService_ConsumeServer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &OmegaServer{
				UnsafeOmegaServiceServer: tt.fields.UnsafeOmegaServiceServer,
			}
			if err := s.Consume(tt.args.stream); (err != nil) != tt.wantErr {
				t.Errorf("Consume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOmegaServer_Pub(t *testing.T) {
	type args struct {
		ctx context.Context
		req *omega.PubReq
	}
	tests := []struct {
		name    string
		args    args
		want    *omega.Response
		wantErr bool
	}{
		{"test", args{context.Background(), &omega.PubReq{
			Topic:     "test_topic",
			HashCode:  0,
			Data:      []byte("test"),
			Priority:  0,
			DelayType: 0,
			DelayTime: 0,
		}}, &omega.Response{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &OmegaServer{}
			got, err := s.Pub(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pub() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOmegaServer_Sub(t *testing.T) {
	type fields struct {
		UnsafeOmegaServiceServer omega.UnsafeOmegaServiceServer
	}
	type args struct {
		ctx context.Context
		req *omega.SubReq
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *omega.SubRsp
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &OmegaServer{
				UnsafeOmegaServiceServer: tt.fields.UnsafeOmegaServiceServer,
			}
			got, err := s.Sub(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sub() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOmegaServer_UpdateSubCfg(t *testing.T) {
	type fields struct {
		UnsafeOmegaServiceServer omega.UnsafeOmegaServiceServer
	}
	type args struct {
		ctx context.Context
		req *omega.SubCfg
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *omega.Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &OmegaServer{
				UnsafeOmegaServiceServer: tt.fields.UnsafeOmegaServiceServer,
			}
			got, err := s.UpdateSubCfg(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateSubCfg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateSubCfg() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateRandomHashCode(t *testing.T) {
	tests := []struct {
		name string
		want uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateRandomHashCode(); got != tt.want {
				t.Errorf("generateRandomHashCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
