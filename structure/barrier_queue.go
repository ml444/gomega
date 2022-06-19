package structure

import "github.com/ml444/scheduler/backend"

//
// 串行模式下的先入先出队列
//
type RetryItem struct {
	item      *backend.Item
	retryAtMs int64
}



