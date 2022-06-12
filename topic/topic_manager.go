package topic

import (
	"errors"
	"fmt"
	log "github.com/ml444/glog"
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

const (
	topicConfigKeyPrefix = "smq_topic_"
)

func (*Manager) getTopicConfigKey(name string) string {
	return fmt.Sprintf("%s%s", topicConfigKeyPrefix, name)
}

func (p *Manager) getTopic(name string) *Topic {
	p.mu.RLock()
	t := p.m[name]
	p.mu.RUnlock()
	return t
}

func (p *Manager) setTopic(topic *Topic, updateChannel bool, updateStore bool) (*Topic, error) {
	var t *Topic
	name := topic.Name
	p.mu.RLock()
	t = p.m[name]
	p.mu.RUnlock()
	up := false
	if t == nil {
		// not found topic, create
		if err := p.checkNaming(name); err != nil {
			return nil, err
		}
		logic := func() error {
			p.mu.Lock()
			defer p.mu.Unlock()
			t = p.m[name]
			if t == nil {
				var err error
				t, err = NewTopic(topic)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

				err = t.init()
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}
				for _, c := range t.subscribers {
					err = c.init()
					if err != nil {
						log.Errorf("err:%v", err)
						return err
					}
				}
				err = p.watchTopic(name)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}
				p.m[name] = t
				up = true
			}
			return nil
		}
		err := logic()
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
	} else {
		// 以etcd实时取出来的NodeGroup为准
		oldTopic, err := topicMgr.getTopicConfig(topic.Name)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		t.group = oldTopic.NodeGroup
		t.withFileGroup.group = oldTopic.NodeGroup

		err = t.init()
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
		for _, v := range t.subscribers {
			v.group = oldTopic.NodeGroup
			err = v.init()
			if err != nil {
				log.Errorf("err:%v", err)
				return nil, err
			}
		}

		t.t.UpdateBy = topic.UpdateBy
		if setTopicCheckUpdate(t.t, topic) {
			t.t.UpdatedAt = utils.Now()
			up = true
			t.subCfg = topic.SubConfig
		}
	}
	if up && updateStore {
		log.Infof("update topic config to etcd")
		err := topicMgr.updateTopicConfig(t.t)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
	} else {
		log.Infof("skip update topic config to etcd")
	}
	if updateChannel {
		t.subscribersMu.Lock()
		defer t.subscribersMu.Unlock()
		newMap := map[string]*smq.Channel{}
		for _, v := range topic.ChannelList {
			newMap[v.Name] = v
		}
		// 新增
		for _, v := range topic.ChannelList {
			if _, ok := t.subscribers[v.Name]; !ok {
				c := NewSubscriber(t.group, v.Name, v)
				err := c.init()
				if err != nil {
					log.Errorf("err:%v", err)
					return nil, err
				}
				t.subscribers[v.Name] = c
			}
		}
		// 删除
		var delKeys []string
		for k, v := range t.subscribers {
			if newMap[k] == nil {
				delKeys = append(delKeys, k)
				v.close()
			}
		}
		for _, v := range delKeys {
			delete(t.subscribers, v)
		}
	}
	return t, nil
}

func (p *Manager) watchTopic(name string) error {
	var err error
	cfgKey := p.getTopicConfigKey(name)
	watcher := config.NewItemWatcher(cfgKey)
	err = watcher.Start(func(ev int, item *config.Item) {
		if (ev == config.ItemCreate || ev == config.ItemUpdate) && item.Val != "" {
			var x smq.Topic
			err = utils.Json2Pb(item.Val, &x)
			if err == nil {
				_, err = p.setTopic(&x, true, false)
				if err != nil {
					log.Errorf("err:%v", err)
				}
			}
		}
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}
func (p *Manager) getOrCreateTopic(name string) (*Topic, error) {
	p.mu.RLock()
	t := p.m[name]
	p.mu.RUnlock()
	if t != nil {
		return t, nil
	}
	if err := p.checkNaming(name); err != nil {
		return nil, err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	t = p.m[name]
	if t != nil {
		return t, nil
	}
	c := config.NewConfig()
	defer func() {
		_ = c.Close()
	}()
	cfgKey := p.getTopicConfigKey(name)
	item, err := c.Get(cfgKey, nil)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	if item == nil || item.Val == "" {
		st := &smq.Topic{
			CreatedAt:   utils.Now(),
			Name:        name,
			ChannelList: nil,
		}
		t, err = NewTopic(st)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		j, err := utils.Pb2Json(st)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		err = c.Set(cfgKey, j)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
	} else {
		var x smq.Topic
		err = utils.Json2Pb(item.Val, &x)
		if err != nil {
			log.Errorf("parse json err:%v, json %s", err, item.Val)
			return nil, err
		}

		t, err = NewTopic(&x)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
	}
	log.Infof("init topic %s", name)
	err = t.init()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	for _, c := range t.subscribers {
		err = c.init()
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
	}
	// log.Infof("watch topic %s", name)
	err = p.watchTopic(name)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	p.m[name] = t
	// log.Infof("added topic %s", name)
	return t, nil
}

// 2021-09-02 不支持修改NodeGroup
func (p *Manager) updateTopicConfig(t *smq.Topic) error {
	x, err := p.getTopicConfig(t.Name)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	t.NodeGroup = x.NodeGroup
	t.UpdateBy = rpc.BindIp

	j, err := utils.Pb2Json(t)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	cfgKey := p.getTopicConfigKey(t.Name)
	c := config.NewConfig()
	defer func() {
		_ = c.Close()
	}()

	err = c.Set(cfgKey, j)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (p *Manager) getTopicConfig(name string) (*smq.Topic, error) {
	c := config.NewConfig()
	defer func() {
		_ = c.Close()
	}()
	cfgKey := p.getTopicConfigKey(name)

	var fromLocalFs bool
	item, err := c.Get(cfgKey, &fromLocalFs)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	var x smq.Topic

	if item.Val == "" {
		return &x, nil
	}

	err = utils.Json2Pb(item.Val, &x)
	if err != nil {
		log.Errorf("parse json err:%v, json %s", err, item.Val)
		return nil, err
	}

	return &x, nil
}

func (p *Manager) delChannel(topicName, channelName string) error {
	topic := p.getTopic(topicName)
	if topic == nil {
		return rpc.CreateError(smq.ErrTopicNotFound)
	}
	if topic.getSubscriber(channelName) == nil {
		return nil
	}
	topic.subscribersMu.Lock()
	defer topic.subscribersMu.Unlock()
	var c *Channel
	if c = topic.subscribers[channelName]; c == nil {
		return nil
	}
	// update config
	var newList []*smq.Channel
	for _, v := range topic.t.ChannelList {
		if v.Name != channelName {
			newList = append(newList, v)
		}
	}
	topic.t.ChannelList = newList
	err := p.updateTopicConfig(topic.t)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	c.close()
	delete(topic.subscribers, channelName)
	return nil
}
func setSubConfig(o *Config, n *smq.SubConfig) bool {
	upCnt := 0
	if o.ServiceName != n.ServiceName {
		o.ServiceName = n.ServiceName
		upCnt++
	}
	if o.ServicePath != n.ServicePath {
		o.ServicePath = n.ServicePath
		upCnt++
	}
	if o.ConcurrentCount != n.ConcurrentCount {
		o.ConcurrentCount = n.ConcurrentCount
		upCnt++
	}
	if o.MaxRetryCount != n.MaxRetryCount {
		o.MaxRetryCount = n.MaxRetryCount
		upCnt++
	}
	if o.MaxExecTimeSeconds != n.MaxExecTimeSeconds {
		o.MaxExecTimeSeconds = n.MaxExecTimeSeconds
		upCnt++
	}
	if o.MaxInQueueTimeSeconds != n.MaxInQueueTimeSeconds {
		o.MaxInQueueTimeSeconds = n.MaxInQueueTimeSeconds
		upCnt++
	}
	if o.BarrierMode != n.BarrierMode {
		o.BarrierMode = n.BarrierMode
		upCnt++
	}
	if o.BarrierCount != n.BarrierCount {
		o.BarrierCount = n.BarrierCount
		upCnt++
	}
	return upCnt > 0
}
func setTopicCheckUpdate(o *smq.Topic, n *smq.Topic) bool {
	upCnt := 0
	// if o.NodeGroup != n.NodeGroup {
	//	o.NodeGroup = n.NodeGroup
	//	upCnt++
	// }
	if o.MaxMsgSize != n.MaxMsgSize {
		o.MaxMsgSize = n.MaxMsgSize
		upCnt++
	}
	if o.SkipConsumeTopic != n.SkipConsumeTopic {
		o.SkipConsumeTopic = n.SkipConsumeTopic
		upCnt++
	}
	if n.SubConfig != nil {
		if o.SubConfig == nil {
			o.SubConfig = &smq.SubConfig{}
		}
		if setSubConfig(o.SubConfig, n.SubConfig) {
			upCnt++
		}
	}
	return upCnt > 0
}
func setChannelCheckUpdate(o *smq.Channel, n *smq.Channel) bool {
	upCnt := 0
	if n.SubConfig != nil {
		if o.SubConfig == nil {
			o.SubConfig = &smq.SubConfig{}
		}
		if setSubConfig(o.SubConfig, n.SubConfig) {
			upCnt++
		}
	}
	return upCnt > 0
}
func (p *Manager) setChannel(topicName string, channel *smq.Channel, updateStore bool) (*Topic, *Channel, error) {
	topic := p.getTopic(topicName)
	if topic == nil {
		return nil, nil, rpc.CreateError(smq.ErrTopicNotFound)
	}
	channelName := channel.Name
	var c *Channel
	up := false
	if c = topic.getSubscriber(channelName); c == nil {
		// 新增
		if err := p.checkNaming(channelName); err != nil {
			return topic, nil, err
		}
		logic := func() error {
			topic.subscribersMu.Lock()
			defer topic.subscribersMu.Unlock()
			// double check
			c = topic.subscribers[channelName]
			if c == nil {
				channel.CreatedAt = utils.Now()
				topic.t.ChannelList = append(topic.t.ChannelList, channel)
				c = NewSubscriber(topic.t.NodeGroup, topicName, channel)
				err := c.init()
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}
				topic.subscribers[channelName] = c
				up = true
			}
			return nil
		}
		err := logic()
		if err != nil {
			log.Errorf("err:%v", err)
			return topic, nil, err
		}
	} else {
		// 修改
		if setChannelCheckUpdate(c.c, channel) {
			c.c.UpdatedAt = utils.Now()
			up = true
			c.subCfg = channel.SubConfig
		}
	}
	if up && updateStore {
		log.Infof("update topic config to etcd - channel")
		err := p.updateTopicConfig(topic.t)
		if err != nil {
			log.Errorf("err:%v", err)
			return topic, nil, err
		}
	} else {
		log.Infof("skip update topic config to etcd - channel")
	}
	return topic, c, nil
}
func (p *Manager) restoreTopics() error {
	c := config.NewConfig()
	defer func() {
		_ = c.Close()
	}()
	items, err := c.ListByPrefix(topicConfigKeyPrefix, 0)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	log.Infof("got %d items", len(items))
	for _, v := range items {
		if v.Val == "" {
			continue
		}
		var x smq.Topic
		err = utils.Json2Pb(v.Val, &x)
		if err != nil {
			log.Errorf("parse json err:%v, json %s", err, v.Val)
			continue
		}
		log.Infof("init topic %s...", x.Name)
		_, err = p.getOrCreateTopic(x.Name)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}
	log.Infof("restore topics success")
	return nil
}
