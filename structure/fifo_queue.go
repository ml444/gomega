package structure

import (
	"encoding/binary"
	"errors"
	"fmt"
	log "github.com/ml444/glog"
	"github.com/ml444/scheduler/backend"
	"go.uber.org/atomic"
	"io"
	"os"
	"sync"
	"time"
)

type BarrierQueueReaderConfig struct {
	MaxCacheIndexFileCount int
}

func NewDefaultBarrierQueueReaderConfig() *BarrierQueueReaderConfig {
	return &BarrierQueueReaderConfig{
		MaxCacheIndexFileCount: 15,
	}
}

type BarrierQueueReaderDelay struct {
	enableDelay         bool
	hash2LastFinishTs   map[uint64]int64
	hash2LastFinishTsMu sync.RWMutex
}
type SerialQueueReader struct {
	backend.fileBase
	cfg          *BarrierQueueReaderConfig
	itemCount    uint32
	readCursor   uint32
	finishFile   *os.File
	finishMap    map[uint32]bool
	recountIndex bool
	indexes      *IndexCache
	popCnt       int
	finishedCnt  atomic.Int32
	delay        BarrierQueueReaderDelay

	indexBuf []byte

	openAt int64
}

func NewBarrierQueueReader(name, dataPath string, seq uint64, cfg *BarrierQueueReaderConfig) *SerialQueueReader {
	if cfg == nil {
		cfg = NewDefaultBarrierQueueReaderConfig()
	}
	return &SerialQueueReader{
		fileBase: backend.fileBase{
			name:     name,
			dataPath: dataPath,
			seq:      seq,
		},
		cfg:       cfg,
		finishMap: map[uint32]bool{},
		indexes:   NewIndexCache(),
	}
}

func (p *SerialQueueReader) statItemCount() error {
	idxPath := p.indexFilePath()
	idxInfo, err := os.Stat(idxPath)
	if err != nil {
		return err
	}
	idxSize := int(idxInfo.Size())
	c := uint32(idxSize / backend.indexItemSize)
	if c < p.itemCount {
		return fmt.Errorf("index file truncated, origin %d, cur %d", p.itemCount, c)
	}
	p.itemCount = c
	return nil
}

func (p *SerialQueueReader) dump() {
	log.Infof("%s.%d: pop %d fin cnt %d total %d sk %d pending %d finished %v",
		p.name, p.seq, p.popCnt, p.finishedCnt,
		//p.indexes.itemTotalCnt, p.indexes.minIndexList.Len(),
		p.indexes.itemTotalCnt-p.popCnt, p.isFinished())
}

func isFileNotFoundError(err error) bool {
	return os.IsNotExist(err)
}

func (p *SerialQueueReader) Init() error {
	var err error
	idxPath := p.indexFilePath()
	datPath := p.dataFilePath()
	finishPath := p.finishFilePath()
	p.indexFile, err = os.OpenFile(idxPath, os.O_RDONLY, 0666)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	p.dataFile, err = os.OpenFile(datPath, os.O_RDONLY, 0666)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	p.finishFile, err = os.OpenFile(finishPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	err = p.statItemCount()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	err = p.loadFinishFile()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	err = p.fillIndexCache()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	p.openAt = time.Now().UnixMilli()
	return nil
}

func (p *SerialQueueReader) close() {
	log.Infof("%s.%d: close", p.name, p.seq)

	if p.indexFile != nil {
		_ = p.indexFile.Close()
		p.indexFile = nil
	}
	if p.dataFile != nil {
		_ = p.dataFile.Close()
		p.dataFile = nil
	}
	if p.finishFile != nil {
		_ = p.finishFile.Close()
		p.finishFile = nil
	}
}


func (p *SerialQueueReader) fillIndexCache() error {
	f := p.indexFile
	if f == nil {
		return errors.New("file not open")
	}
	if p.dataCorruption {
		return errors.New("data corruption")
	}
	if p.readCursor < p.itemCount {
		// 跳过已处理的
		for p.readCursor < p.itemCount {
			if p.finishMap[p.readCursor] {
				p.readCursor++
			} else {
				break
			}
		}
		if int(p.itemCount-p.readCursor) <= 0 {
			log.Infof("skip load index, all finished")
			return nil
		}

		// init buffer
		if len(p.indexBuf) == 0 {
			log.Infof("alloc index buf for %s.%d", p.name, p.seq)
			p.indexBuf = make([]byte, (1024*1024)/backend.indexItemSize*backend.indexItemSize)
		}

		eof := false
		for p.readCursor < p.itemCount && !eof {
			n, err := f.ReadAt(p.indexBuf, int64(p.readCursor)*backend.indexItemSize)
			if err != nil {
				if err == io.EOF {
					eof = true
				} else {
					log.Errorf("err:%v", err)
					return err
				}
			}
			log.Infof("%s seq %d: load index at %d with size %d, size %d ret",
				p.name, p.seq, p.readCursor*backend.indexItemSize, len(p.indexBuf), n)
			if n <= 0 {
				break
			}
			n = n / backend.indexItemSize
			b := binary.LittleEndian
			var items []*backend.Item
			for i := 0; i < n; i++ {
				index := p.readCursor
				p.readCursor++
				if p.finishMap[index] {
					continue
				}
				ptr := p.indexBuf[i*backend.indexItemSize:]
				var begMarker, endMarker uint16
				begMarker = b.Uint16(ptr)
				ptr = ptr[2:]
				endMarker = b.Uint16(ptr[28:])
				if begMarker != backend.itemBegin {
					p.dataCorruption = true
					return errors.New("invalid index item begin marker")
				}
				if endMarker != backend.itemEnd {
					p.dataCorruption = true
					return errors.New("invalid index item end marker")
				}
				var idx backend.Item
				idx.CreatedAt = b.Uint64(ptr)
				//idx.CorpId = b.Uint32(ptr[4:])
				//idx.AppId = b.Uint32(ptr[8:])
				idx.Hash = b.Uint64(ptr[8:])
				idx.DelayType, idx.DelayValue, idx.Priority = backend.unpackMisc(b.Uint32(ptr[16:]))
				idx.offset = b.Uint32(ptr[20:])
				idx.size = b.Uint32(ptr[24:])
				idx.Index = index
				idx.Seq = p.seq
				if idx.size < backend.dataItemExtraSize {
					p.dataCorruption = true
					return fmt.Errorf("invalid data size %d, min than data item extra size", idx.size)
				}
				if idx.DelayType == backend.DelayTypeRelate {
					p.delay.enableDelay = true
				}
				items = append(items, &idx)
			}
			if len(items) > 0 {
				p.indexes.addItems(items)
			}
		}
	}
	return nil
}

func (p *SerialQueueReader) loadFinishFile() error {
	if p.finishFile == nil {
		panic("Unreachable")
	}
	const (
		loadBufSize = 10240 * 4
	)
	buf := make([]byte, loadBufSize)
	_, err := p.finishFile.Seek(0, io.SeekStart)
	if err != nil {
		log.Errorf("seek err:%v", err)
		return err
	}
	for {
		n, err := p.finishFile.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Errorf("err:%v", err)
			return err
		}
		if n > 0 {
			if n%4 != 0 {
				log.Warnf("invalid size of finish file read return %d", n)
				p.dataCorruption = true
				return errors.New("invalid size of finish file")
			}
			b := binary.LittleEndian
			for i := 0; i < n; i += 4 {
				idx := b.Uint32(buf[i:])
				p.finishMap[idx] = true
			}
		}
	}
	log.Infof("finished len %d", len(p.finishMap))
	return nil
}

func (p *SerialQueueReader) readData(item *backend.Item) error {
	d := p.dataFile
	if d == nil {
		return errors.New("file not opened")
	}
	dataBuf := make([]byte, item.size)
	var read uint32
	for read < item.size {
		n, err := d.ReadAt(dataBuf[read:], int64(item.offset)+int64(read))
		if err != nil {
			if err == io.EOF {
				return errors.New("data file truncated")
			}
			log.Errorf("err:%v", err)
			return err
		} else if n > 0 {
			read += uint32(n)
		} else {
			return errors.New("data file truncated")
		}
	}
	ptr := dataBuf[:]
	b := binary.LittleEndian
	begMarker := b.Uint16(ptr)
	item.Data = ptr[2 : item.size-2]
	endMarker := b.Uint16(ptr[item.size-2:])
	if begMarker != backend.itemBegin {
		p.dataCorruption = true
		return errors.New("invalid data begin marker")
	}
	if endMarker != backend.itemEnd {
		p.dataCorruption = true
		return errors.New("invalid data end marker")
	}
	return nil
}

func (p *SerialQueueReader) retry(item *backend.Item, delayMs uint32, wm **WarnMsg) {
	if p.popCnt <= 0 {
		*wm = &WarnMsg{
			Label: fmt.Sprintf("%s.%d: invalid pop cnt", p.name, p.seq),
		}
		return
	}
	p.popCnt--
	p.indexes.retry(item, delayMs)
}

func (p *SerialQueueReader) popSkipHash(skipHash map[uint64]int, barrierCount int, checkCount bool, wm **WarnMsg) (*backend.Item, error) {
	if p.recountIndex {
		p.recountIndex = false

		err := p.statItemCount()
		if err != nil {
			log.Errorf("err:%v", err)
			p.recountIndex = true
			return nil, err
		}
		err = p.fillIndexCache()
		if err != nil {
			log.Errorf("err:%v", err)
			p.recountIndex = true
			return nil, err
		}
	}

	item := p.indexes.popSkipHash(skipHash, barrierCount, checkCount, &p.delay, wm)
	if item == nil {
		return nil, nil
	}
	p.popCnt++

	err := p.readData(item)
	if err != nil {
		log.Errorf("err:%v", err)
		// drop msg
		p.finishedCnt.Add(1)
		return nil, err
	}
	return item, nil
}

func (p *SerialQueueReader) isFinished() bool {
	if len(p.indexes.hashList) == 0 {
		res := (p.popCnt <= int(p.finishedCnt.Load())) && !p.recountIndex
		if res {
			return true
		}
	}
	return false
}

func (p *SerialQueueReader) getQueueLen() int64 {
	var l int64
	t := p.indexes.itemTotalCnt
	pc := p.popCnt
	if t > pc {
		l = int64(t - pc)
	}
	return l
}
