package structure

import (
	"github.com/huandu/skiplist"
	log "github.com/ml444/glog"
	"github.com/ml444/scheduler/backend"
	"time"
)

const MaxPriority = 3
const PriorityMask = 0x3

type PriorityList struct {
	ll [MaxPriority][]*backend.Item
}

type IndexCacheItemList struct {
	list      []*backend.Item
	pl        *PriorityList
	retryList *MinHeap
}

type IndexCache struct {
	minIndexList *skiplist.SkipList
	hashList     map[uint64]*IndexCacheItemList

	itemTotalCnt int
}

func NewIndexCache() *IndexCache {
	return &IndexCache{
		minIndexList: skiplist.New(skiplist.Uint32),
		hashList:     map[uint64]*IndexCacheItemList{},
	}
}

func (p *IndexCacheItemList) appendItem(item *backend.Item) {
	if item.Priority == 0 {
		p.list = append(p.list, item)
	} else {
		pr := (item.Priority & PriorityMask) - 1
		if p.pl == nil {
			p.pl = &PriorityList{}
		}
		p.pl.ll[pr] = append(p.pl.ll[pr], item)
	}
}

func (p *IndexCacheItemList) popItem(check func(item *backend.Item) bool) *backend.Item {
	if p.pl != nil {
		for i := MaxPriority - 1; i >= 0; i-- {
			x := p.pl.ll[i]
			if len(x) > 0 {
				f := x[0]
				if check == nil || check(f) {
					x = x[1:]
					p.pl.ll[i] = x
					return f
				} else {
					// 严格先跑高优的
					return nil
				}
			}
		}
	}

	if len(p.list) > 0 {
		f := p.list[0]
		if check == nil || check(f) {
			p.list = p.list[1:]
			return f
		}
	}

	return nil
}

func (p *IndexCacheItemList) empty() bool {
	if p.pl != nil {
		for i := MaxPriority - 1; i >= 0; i-- {
			x := p.pl.ll[i]
			if len(x) > 0 {
				return false
			}
		}
	}
	return len(p.list) == 0
}

func (p *IndexCacheItemList) peekItem(check func(item *backend.Item) bool) *backend.Item {
	if p.pl != nil {
		for i := MaxPriority - 1; i >= 0; i-- {
			x := p.pl.ll[i]
			if len(x) > 0 {
				f := x[0]
				if check == nil || check(f) {
					return f
				}
			}
		}
	}

	if len(p.list) > 0 {
		f := p.list[0]
		if check == nil || check(f) {
			return f
		}
	}

	return nil
}

func (p *IndexCache) addItems(items []*backend.Item) {
	m := p.minIndexList
	h := p.hashList
	for _, v := range items {
		x := h[v.Hash]
		if x == nil {
			x = &IndexCacheItemList{}
			m.Set(v.Index, v.Hash)
		}
		x.appendItem(v)
		h[v.Hash] = x
	}
	p.itemTotalCnt += len(items)
}

func (p *IndexCache) retry(item *backend.Item, delayMs uint32) {
	m := p.minIndexList
	h := p.hashList
	x := h[item.Hash]
	if x == nil {
		x = &IndexCacheItemList{}
		h[item.Hash] = x
		m.Set(item.Index, item.Hash)
	}
	if x.retryList == nil {
		x.retryList = NewMinHeap()
	}
	retryAtMs := time.Now().UnixMilli() + int64(delayMs)
	x.retryList.PushEl(&MinHeapElement{
		Value: &RetryItem{
			item:      item,
			retryAtMs: retryAtMs,
		},
		Priority: retryAtMs,
	})
}

func (p *IndexCache) popSkipHash(skipHash map[uint64]int, barrierCount int, checkCount bool, delay *BarrierQueueReaderDelay, wm **WarnMsg) *backend.Item {
	m := p.minIndexList
	h := p.hashList
	c := m.Front()
	var nowMs int64
	for ; c != nil; c = c.Next() {
		hash := c.Value.(uint64)
		if checkCount && skipHash != nil {
			count := skipHash[hash]
			if count >= barrierCount {
				continue
			}
		}
		x := h[hash]
		if x == nil {
			*wm = &WarnMsg{
				Label: "pop len(list) == 0",
			}
			return nil
		}
		// 看看 retry list
		var item *backend.Item
		if x.retryList != nil {
			if nowMs == 0 {
				nowMs = time.Now().UnixMilli()
			}
			top := x.retryList.PeekEl()
			if top.Priority > nowMs {
				continue
			}
			item = x.retryList.PopEl().Value.(*RetryItem).item
			if x.retryList.Len() == 0 {
				x.retryList = nil
			}
		} else {
			item = x.popItem(func(item *backend.Item) bool {
				if delay.enableDelay && item.DelayType == backend.DelayTypeRelate {
					if nowMs == 0 {
						nowMs = time.Now().UnixMilli()
					}
					var last int64
					if delay.hash2LastFinishTs != nil {
						delay.hash2LastFinishTsMu.RLock()
						last = delay.hash2LastFinishTs[item.Hash]
						delay.hash2LastFinishTsMu.RUnlock()
					}
					if last+int64(item.DelayValue) > nowMs {
						return false
					}
				}
				return true
			})
		}
		if item != nil {
			m.Remove(c.Key())
			if x.empty() && x.retryList == nil {
				delete(h, hash)
			} else {
				if x.retryList != nil {
					top := x.retryList.PeekEl()
					m.Set(top.Value.(*RetryItem).item.Index, hash)
				} else {
					f := x.peekItem(nil)
					if f != nil {
						m.Set(f.Index, hash)
					}
				}
			}
		}
		if item != nil {
			return item
		}
	}
	return nil
}
func (p *IndexCache) dump() {
	log.Infof("skip list len %d, total %d", p.minIndexList.Len(), p.itemTotalCnt)
}
