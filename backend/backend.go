package backend

import "github.com/ml444/scheduler/publish"

type IBackendReader interface {
	Read()
}

type IBackendWriter interface {
	Write(req *publish.PubReq, data []byte)
}