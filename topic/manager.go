package topic

import (
	"errors"
	"fmt"
	"sync"
)

type Manager struct {
	m  map[string]*Topic
	mu sync.RWMutex
}

var topicMgr *Manager

func (*Manager) checkNaming(name string) error {
	l := len(name)
	if l == 0 {
		return errors.New("name is empty")
	}
	if l > 64 {
		return fmt.Errorf( "the length of the name[%s] is greater than 64", name)
	}
	for i := 0; i < l; i++ {
		if name[i] == '.' || name[i] == '/' {
			return fmt.Errorf("name[%s] has invalid char", name)
		}
	}
	return nil
}

func TODO() {
	// 根据Topic本身的优先级，选择一个合适的服务节点，更新etcd
	// sdk：topic
}