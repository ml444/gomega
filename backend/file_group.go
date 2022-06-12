package backend

import (
	"errors"
	"fmt"
	log "github.com/ml444/glog"
	"github.com/ml444/scheduler/mfile"
	"github.com/ml444/scheduler/publish"
	"github.com/petermattis/goid"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type FileGroup struct {
	name        string
	logName     string
	fg          *mfile.FileGroup
	fgMu        sync.Mutex
	qr          mfile.IQueueGroupReader
	qrMu        sync.Mutex
	subCfg      *smq.SubConfig
	topicName   string
	channelName string
	hasInit     bool
	group       string
}
func (p *FileGroup) closeFileGroup() {
	if p.fg != nil {
		p.fgMu.Lock()
		defer p.fgMu.Unlock()
		if p.fg != nil {
			p.fg.Close()
			p.fg = nil
		}
	}
}

func (p *FileGroup) closeReader() {
	if p.qr != nil {
		p.qrMu.Lock()
		defer p.qrMu.Unlock()
		if p.qr != nil {
			p.qr.Close()
			p.qr = nil
		}
	}
}

func (p *FileGroup) close() {
	p.closeFileGroup()
	p.closeReader()
}

func (p *FileGroup) initFileGroup() error {
	if p.fg != nil {
		return nil
	}
	p.fgMu.Lock()
	defer p.fgMu.Unlock()
	if p.fg != nil {
		return nil
	}
	p.fg = mfile.NewFileGroup(p.name, s.Conf.DataPath, nil)
	err := p.fg.Init()
	if err != nil {
		return err
	}
	return nil
}
const (
	defaultConsumeConcurrentCount    = 100
	defaultConsumeMaxExecTimeSeconds = 60
	defaultConsumeMaxRetryCount      = 5
)




func (p *FileGroup) initReader() error {
	if p.qr != nil {
		return nil
	}
	if p.subCfg != nil && p.subCfg.BarrierMode {
		return p.initReaderBarrierMode()
	}
	p.qrMu.Lock()
	defer p.qrMu.Unlock()
	if p.qr != nil {
		return nil
	}
	getConcurrentCount := func() uint32 {
		s := p.subCfg
		if s == nil || s.ConcurrentCount == 0 {
			return defaultConsumeConcurrentCount
		}
		return s.ConcurrentCount
	}
	qgr := mfile.NewQueueGroupReader(p.name, s.Conf.DataPath, nil, 8000)
	qgr.ReportQueueLenFunc = func(ql int64) {
		reportQueueLen(p.topicName, p.channelName, ql)
	}
	p.qr = qgr
	err := p.qr.Init()
	if err != nil {
		log.Errorf("err:%v", err)
		reportWarnMsg(p.topicName, fmt.Sprintf("init reader err %v", err))
		return err
	}
	msgChan := qgr.GetMsgChan()
	errChan := p.qr.GetErrChan()
	warnChan := p.qr.GetWarnChan()
	closeChan := p.qr.GetCloseChan()
	go func() {
		concurrent := getConcurrentCount()
		tk := time.NewTicker(5 * time.Second)
		defer tk.Stop()
		log.Infof("%s: sub cfg %+v", p.logName, p.subCfg)
		var workers []*consumeWorker
		var wg sync.WaitGroup
		startWorkers := func(retryList []*retryItem) {
			// 把 retry 平分一下
			retryCntPerRoutine := (len(retryList) + int(concurrent) - 1) / int(concurrent)
			assignIdx := 0
			for i := 0; i < int(concurrent); i++ {
				w := newConsumeWorker(&wg, i, msgChan, p)
				workers = append(workers, w)
				for j := 0; j < retryCntPerRoutine; j++ {
					if assignIdx < len(retryList) {
						r := retryList[assignIdx]
						assignIdx++
						w.retryList.PushEl(&utils.MinHeapElement{
							Value:    r,
							Priority: int64(r.nextExecAt),
						})
					} else {
						break
					}
				}
				wg.Add(1)
				go w.start()
			}
			if assignIdx < len(retryList) {
				panic(fmt.Sprintf("bug: assign idx %d, retry list len %d", assignIdx, len(retryList)))
			}
		}
		stopWorkers := func() (retryList []*retryItem) {
			for _, w := range workers {
				w.notifyExit()
			}
			wg.Wait()
			// 收集回所有的旧任务
			for _, w := range workers {
				if w.retryList != nil {
					for {
						t := w.retryList.PopEl()
						if t == nil {
							break
						}
						retryList = append(retryList, t.Value.(*retryItem))
					}
				}
			}
			workers = nil
			return
		}
		startWorkers(nil)
		log.Infof("%s: reader loop", p.logName)
		for {
			select {
			case err = <-errChan:
				reportWarnMsg(p.topicName, fmt.Sprintf("mq err: %v", err.Error()))
			case wm := <-warnChan:
				reportWarnMsg(p.topicName, wm.Label)
			case <-closeChan:
				goto OUT
			case <-tk.C:
				c := getConcurrentCount()
				if c != concurrent {
					log.Infof("%s: change concurrent %d to %d", p.logName, concurrent, c)
					concurrent = c
					// 调整并发数
					log.Debugf("%s: stopping", p.logName)
					retryList := stopWorkers()
					log.Debugf("%s: starting", p.logName)
					startWorkers(retryList)
					log.Debugf("%s: change concurrent finished", p.logName)
				}
			}
		}
	OUT:
		stopWorkers()
	}()
	return nil
}
func (p *FileGroup) initReaderBarrierMode() error {
	if p.qr != nil {
		return nil
	}
	p.qrMu.Lock()
	defer p.qrMu.Unlock()
	if p.qr != nil {
		return nil
	}
	getConcurrentCount := func() uint32 {
		s := p.subCfg
		if s == nil || s.ConcurrentCount == 0 {
			return defaultConsumeConcurrentCount
		}
		return s.ConcurrentCount
	}
	p.qr = mfile.NewBarrierQueueGroupReader(p.name, s.Conf.DataPath, nil)
	err := p.qr.Init()
	if err != nil {
		log.Errorf("err:%v", err)
		warning.ReportMsg(nil, fmt.Sprintf("smq: %s: init reader err:%v", p.logName, err))
		return err
	}
	errChan := p.qr.GetErrChan()
	warnChan := p.qr.GetWarnChan()
	closeChan := p.qr.GetCloseChan()
	go func() {
		concurrent := getConcurrentCount()
		tk := time.NewTicker(5 * time.Second)
		defer tk.Stop()
		log.Infof("%s: sub cfg %+v", p.logName, p.subCfg)
		consumer := NewBarrierConsumer(p, int(concurrent), nil)
		consumer.start()
		log.Infof("%s: reader loop", p.logName)
		for {
			select {
			case err = <-errChan:
				reportWarnMsg(p.topicName, fmt.Sprintf("mq err: %v", err.Error()))
			case wm := <-warnChan:
				reportWarnMsg(p.topicName, wm.Label)
			case <-closeChan:
				goto OUT
			case <-tk.C:
				c := getConcurrentCount()
				if c != concurrent {
					log.Infof("change concurrent %d to %d", concurrent, c)
					concurrent = c
					// 调整并发数
					log.Debugf("stopping consumer...")
					consumer.stop()
					old := consumer
					consumer = NewBarrierConsumer(p, int(concurrent), old)
					log.Debugf("starting consumer...")
					consumer.start()
				}

				reportBarrierQueueReaderNum(p.topicName, p.channelName, p.qr.GetReaderNum())
			}
		}
	OUT:
		consumer.stop()
	}()
	return nil
}

func (p *FileGroup) groupOk() bool {
	if p.group != s.nodeCfg.Group {
		if s.nodeCfg.IgnoreGroupInDevRole && rpc.Meta.IsDevRole() {
			return true
		} else if s.nodeCfg.IgnoreGroupInTestRole && rpc.Meta.IsTestRole() {
			return true
		} else {
			return false
		}
	}
	return true
}

func (p *FileGroup) init() error {
	if p.hasInit {
		return nil
	}
	if !p.groupOk() {
		log.Warnf("node group not match, my group is %s, topic group %s, skip init",
			s.nodeCfg.Group, p.group)
		return nil
	}
	log.Infof("init file group %s", p.logName)
	err := p.initFileGroup()
	if err != nil {
		return err
	}
	log.Infof("init reader %s", p.logName)
	err = p.initReader()
	if err != nil {
		return err
	}
	p.hasInit = true
	return nil
}
func (p *FileGroup) Write(req *publish.PubReq, data []byte) error {
	if p.fg == nil {
		return errors.New("file group not init")
	}
	if p.qr == nil {
		return errors.New("queue reader not init")
	}
	if !p.groupOk() {
		return rpc.InvalidArg("node group not match")
	}
	it := &mfile.Item{
		Hash:       req.Hash,
		Data:       data,
		DelayType:  req.DelayType,
		DelayValue: req.DelayValue,
		Priority:   uint8(req.Priority),
	}
	err := p.fg.Write(it)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	// 入队量上报
	p.qr.NotifyWrite(it.Seq)
	return nil
}
