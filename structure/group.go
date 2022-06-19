package structure

import (
	"errors"
	"fmt"
	log "github.com/ml444/glog"
	"github.com/ml444/scheduler/backend"
	"sync"
	"time"
)

type BarrierQueueGroupReader struct {
	name                     string
	dataPath                 string
	readers                  []*SerialQueueReader
	minSeq                   uint64
	maxSeq                   uint64
	cfg                      *BarrierQueueReaderConfig
	hasInit                  bool
	initMu                   sync.Mutex
	exitChan                 chan int
	notifyWriteChan          chan uint64
	errChan                  chan error
	warnChan                 chan *WarnMsg
	finishChan               chan *FinishReq
	closeChan                chan int
	notifyWriteChan4Consumer chan uint64
	asyncMsgConfirmChan      chan string
}

func NewBarrierQueueGroupReader(name, dataPath string, cfg *BarrierQueueReaderConfig) *BarrierQueueGroupReader {
	if cfg == nil {
		cfg = NewDefaultBarrierQueueReaderConfig()
	}
	return &BarrierQueueGroupReader{
		name:                     name,
		dataPath:                 dataPath,
		maxSeq:                   0,
		cfg:                      cfg,
		exitChan:                 make(chan int),
		notifyWriteChan:          make(chan uint64, 1000),
		errChan:                  make(chan error, 1000),
		warnChan:                 make(chan *WarnMsg, 100),
		finishChan:               make(chan *FinishReq, 1000),
		closeChan:                make(chan int),
		notifyWriteChan4Consumer: make(chan uint64, 1000),
		asyncMsgConfirmChan:      make(chan string, 1000),
	}
}

func (p *BarrierQueueGroupReader) GetAsyncMsgConfirmChan() chan string {
	return p.asyncMsgConfirmChan
}

func (p *BarrierQueueGroupReader) ConfirmAsyncMsg(msgId string) {
	p.asyncMsgConfirmChan <- msgId
}

func (p *BarrierQueueGroupReader) cleanFinishedReader() {
	if len(p.readers) <= 1 {
		return
	}
	var n []*SerialQueueReader
	closeCnt := 0
	last := p.readers[len(p.readers)-1]
	for i := 0; i < len(p.readers)-1; i++ {
		skipCheckFinished := false

		if i+1 == len(p.readers)-1 {
			// 这个文件表示 max seq - 1，当快速切换文件时， max seq - 1 的文件，
			// 可能写入了文件，但是 delay 了一会会 notify 过来
			// 导致新写入的 msg 没有被消费，就 close 掉了
			if last.openAt+30 > time.Now().UnixMilli() {
				skipCheckFinished = true
			}
		}
		r := p.readers[i]
		if !skipCheckFinished && r.isFinished() {
			r.close()
			closeCnt++
		} else {
			n = append(n, r)
		}
	}
	if closeCnt > 0 {
		n = append(n, p.readers[len(p.readers)-1])
		p.readers = n
	}
}

func (p *BarrierQueueGroupReader) isReaderOpened(seq uint64) bool {
	for _, v := range p.readers {
		if v.seq == seq {
			return true
		}
	}
	return false
}

func (p *BarrierQueueGroupReader) openReader() error {
	p.cleanFinishedReader()

	if p.minSeq == 0 {
		return nil
	}
	for p.minSeq <= p.maxSeq && len(p.readers) < p.cfg.MaxCacheIndexFileCount {
		if !p.isReaderOpened(p.minSeq) {
			reader := NewBarrierQueueReader(p.name, p.dataPath, p.minSeq, p.cfg)
			err := reader.Init()
			if err != nil {
				reader.close()
				reader = nil
				log.Errorf("err:%v", err)
				if !isFileNotFoundError(err) {
					return err
				}
			}
			if reader != nil {
				if reader.isFinished() && p.minSeq < p.maxSeq {
					log.Infof("%s %d has finished, skip", reader.name, reader.seq)
					reader.close()
					reader = nil
				} else {
					log.Infof("init %s %d finished", reader.name, reader.seq)
					reader.dump()
					p.readers = append(p.readers, reader)
				}
			}
		}
		if p.minSeq < p.maxSeq {
			p.minSeq++
		} else {
			break
		}
	}
	return nil
}

func (p *BarrierQueueGroupReader) GetReaderNum() int {
	if p == nil {
		return 0
	}

	return len(p.readers)
}

func (p *BarrierQueueGroupReader) ioLoop() {
	log.Infof("start io loop")

	const (
		writeBufSize = 10240 * 4
	)
	writeBuf := make([]byte, writeBufSize)

	flushFinish := func(r *FinishReq) int {
		p.SetLastFinishTs(r.seq, r.hash)

		type wrap struct {
			seq  uint64
			list []uint32
		}
		all := map[uint64]*wrap{}
		var last *wrap
		for r != nil {
			if last == nil || last.seq != r.seq {
				last = all[r.seq]
				if last == nil {
					last = &wrap{
						seq: r.seq,
					}
					all[r.seq] = last
				}
			}
			last.list = append(last.list, r.index)
			select {
			case r = <-p.finishChan:
				p.SetLastFinishTs(r.seq, r.hash)
			default:
				r = nil
			}
		}
		var fn int
		for _, v := range all {
			var reader *SerialQueueReader
			for _, x := range p.readers {
				if x.seq == v.seq {
					reader = x
					break
				}
			}
			if reader != nil {
				err := batchWriteFinishFile(reader.finishFile, reader.finishMap, v.list, writeBuf)
				if err != nil {
					log.Errorf("err:%v", err)
					p.sendErr(err)
				}
				reader.finishedCnt.Add(int32(len(v.list)))
				if reader.isFinished() {
					fn++
				}
			}
		}
		return fn
	}
	tk := time.NewTicker(time.Minute)
	for {
		select {
		case <-tk.C:
			// 兜底
			err := p.openReader()
			if err != nil {
				log.Errorf("err:%v", err)
				p.sendErr(err)
			}

		case seq := <-p.notifyWriteChan:
			if p.minSeq == 0 {
				p.minSeq = 1
			}
			for _, r := range p.readers {
				if r.seq == seq {
					r.recountIndex = true
					break
				}
			}
			if seq > p.maxSeq {
				p.maxSeq = seq
				err := p.openReader()
				if err != nil {
					log.Errorf("err:%v", err)
					p.sendErr(err)
				}
			}
			select {
			case p.notifyWriteChan4Consumer <- seq:
			default:
			}
		case <-p.exitChan:
			goto EXIT
		case r := <-p.finishChan:
			fn := flushFinish(r)
			if fn > 0 {
				err := p.openReader()
				if err != nil {
					log.Errorf("err:%v", err)
					p.sendErr(err)
				}
			}
		}
	}
EXIT:
	select {
	case p.closeChan <- 1:
	default:
	}
}

func (p *BarrierQueueGroupReader) Init() error {
	p.initMu.Lock()
	defer p.initMu.Unlock()
	if p.hasInit {
		p.warnChan <- &WarnMsg{
			Label: fmt.Sprintf("%s: duplicate init", p.name),
		}
		return errors.New("duplicate init")
	}
	err := backend.scanInitFileGroup(p.name, p.dataPath, &p.minSeq, &p.maxSeq)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	if p.minSeq > p.maxSeq {
		p.warnChan <- &WarnMsg{
			Label: fmt.Sprintf("%s: minSeq > maxSeq", p.name),
		}
		return errors.New("unreachable")
	}
	err = p.openReader()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	go p.ioLoop()
	p.hasInit = true
	return nil
}

func (p *BarrierQueueGroupReader) sendErr(err error) {
	select {
	case p.errChan <- err:
	default:
	}
}

func (p *BarrierQueueGroupReader) SetLastFinishTs(seq uint64, hash uint64) {
	for _, reader := range p.readers {
		if reader.seq == seq {
			if reader.delay.enableDelay {
				reader.delay.hash2LastFinishTsMu.Lock()

				if reader.delay.hash2LastFinishTs == nil {
					reader.delay.hash2LastFinishTs = map[uint64]int64{}
				}
				reader.delay.hash2LastFinishTs[hash] = time.Now().UnixMilli()

				reader.delay.hash2LastFinishTsMu.Unlock()
			}
			break
		}
	}
}

func (p *BarrierQueueGroupReader) SetFinish(item *backend.Item) {
	p.finishChan <- &FinishReq{
		seq:   item.Seq,
		index: item.Index,
		hash:  item.Hash,
	}
}

func (p *BarrierQueueGroupReader) NotifyWrite(seq uint64) {
	select {
	case p.notifyWriteChan <- seq:
	default:
	}
}

func (p *BarrierQueueGroupReader) Close() {
	select {
	case p.exitChan <- 1:
	default:
	}
}

func (p *BarrierQueueGroupReader) popSkipHash(skipHash map[uint64]int, barrierCount int, checkCount bool) (*backend.Item, error) {
	for _, reader := range p.readers {
		var wm *WarnMsg
		item, err := reader.popSkipHash(skipHash, barrierCount, checkCount, &wm)
		if wm != nil {
			wm.Label = fmt.Sprintf("%s.%d: %s", p.name, reader.seq, wm.Label)
			p.warnChan <- wm
		}
		if err != nil {
			log.Errorf("err:%v", err)
			p.sendErr(err)
			return nil, err
		}
		if item != nil {
			return item, nil
		}
	}
	return nil, nil
}

func (p *BarrierQueueGroupReader) PopSkipHash(skipHash map[uint64]int, barrierCount int) (*backend.Item, error) {
	if barrierCount <= 0 {
		barrierCount = 1
	}

	item, err := p.popSkipHash(skipHash, barrierCount, true)
	if err != nil {
		return nil, err
	}
	if item != nil {
		return item, nil
	}

	if barrierCount > 1 {
		return p.popSkipHash(skipHash, barrierCount, false)
	}

	return nil, nil
}

func (p *BarrierQueueGroupReader) Retry(item *backend.Item, delayMs uint32) {
	if delayMs > 3600*1000 {
		log.Warnf("%s retry delay too large %d", p.name, delayMs)

		// 拉群经常报这个，但又是正常的，先屏蔽掉
		//p.warnChan <- &WarnMsg{
		//	Label: "retry delay too large",
		//	Msg:   fmt.Sprintf("%s: retry delay %d ms", p.name, delayMs),
		//}
	}
	item.Data = nil
	for _, reader := range p.readers {
		if reader.seq == item.Seq {
			var wm *WarnMsg
			reader.retry(item, delayMs, &wm)
			if wm != nil {
				p.warnChan <- wm
			}
			return
		}
	}

	log.Warnf("not found reader %s seq %d", p.name, item.Seq)

	p.warnChan <- &WarnMsg{
		Label: fmt.Sprintf("%s: not found reader", p.name),
	}
}

func (p *BarrierQueueGroupReader) GetErrChan() chan error {
	return p.errChan
}

func (p *BarrierQueueGroupReader) GetWarnChan() chan *WarnMsg {
	return p.warnChan
}

func (p *BarrierQueueGroupReader) GetCloseChan() chan int {
	return p.closeChan
}

func (p *BarrierQueueGroupReader) GetQueueLen() int64 {
	var l int64
	for _, r := range p.readers {
		l += r.getQueueLen()
	}
	return l
}

func (p *BarrierQueueGroupReader) GetNotifyWriteChan4Consumer() chan uint64 {
	return p.notifyWriteChan4Consumer
}

func (p *BarrierQueueGroupReader) DumpReaders() {
	for _, v := range p.readers {
		log.Infof("dump reader with seq %d", v.seq)
		v.dump()
		log.Infof("===")
	}
}
