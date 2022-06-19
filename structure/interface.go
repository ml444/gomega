package structure

import "github.com/ml444/scheduler/backend"

type WarnMsg struct {
	Label string
	Msg   string
}

type IQueueGroupReader interface {
	Init() error
	Close()
	NotifyWrite(seq uint64)
	GetErrChan() chan error
	GetCloseChan() chan int
	GetAsyncMsgConfirmChan() chan string
	SetFinish(item *backend.Item)
	SetLastFinishTs(seq uint64, hash uint32)

	GetWarnChan() chan *WarnMsg

	GetQueueLen() int64

	ConfirmAsyncMsg(msgId string)

	GetReaderNum() int
}
