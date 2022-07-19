package brokers

/*
index:
|-begin-|-CreatedAt-|-Sequence-|-HashCode-|-Offset-|-Size-|-end-|
| 2		| 4			| 8		   | 8		  | 4	   | 4	  | 2	|
*/
import (
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	itemBegin = uint16(0x1234)
	itemEnd   = uint16(0x5678)

	indexItemSize     = 32
	dataItemExtraSize = 4

	isLittle = true
)

type Endian struct {
	isLittle bool
}

func getEndian() binary.ByteOrder {
	if isLittle {
		return binary.LittleEndian
	} else {
		return binary.BigEndian
	}
}

type Item struct {
	Sequence  uint64
	HashCode  uint64
	CreatedAt uint32
	Partition uint32
	Offset    uint32
	Size      uint32

	RetryCount uint32
	DelayType  uint32
	DelayValue uint32
	Priority   uint32
	Data       []byte

	//Namespace string
}

func (i *Item) FillIndex(buf []byte) error {
	b := getEndian()
	begMark := b.Uint16(buf)
	endMark := b.Uint16(buf[30:])
	if begMark != itemBegin {
		//p.dataCorruption = true
		return errors.New("invalid index item begin marker")
	}
	if endMark != itemEnd {
		//p.dataCorruption = true
		return errors.New("invalid index item end marker")
	}
	i.CreatedAt = b.Uint32(buf[2:])
	i.Sequence = b.Uint64(buf[6:])
	i.HashCode = b.Uint64(buf[14:])
	i.Offset = b.Uint32(buf[22:])
	i.Size = b.Uint32(buf[26:])
	if i.Size < dataItemExtraSize {
		//p.dataCorruption = true
		return fmt.Errorf("invalid data size %d, min than data item extra size", i.Size)
	}
	return nil
}

func (i *Item) FillData(dataBuf []byte) error {
	ptr := dataBuf[:]
	b := getEndian()
	begMark := b.Uint16(ptr)
	i.Data = ptr[2 : i.Size-2]
	endMark := b.Uint16(ptr[i.Size-2:])
	if begMark != itemBegin {
		//q.dataCorruption = true
		return errors.New("invalid data begin marker")
	}
	if endMark != itemEnd {
		//q.dataCorruption = true
		return errors.New("invalid data end marker")
	}
	return nil
}

func (i *Item) Marshal2Index(buf []byte) {
	b := getEndian()
	b.PutUint16(buf[0:2], itemBegin)
	b.PutUint32(buf[2:6], i.CreatedAt)
	b.PutUint64(buf[6:14], i.Sequence)
	b.PutUint64(buf[14:22], i.HashCode)
	b.PutUint32(buf[22:26], i.Offset)
	b.PutUint32(buf[26:30], i.Size)
	b.PutUint16(buf[30:32], itemEnd)
}

func (i *Item) Marshal2Data(buf []byte) {
	b := getEndian()
	b.PutUint16(buf[0:2], itemBegin)
	copy(buf[2:], i.Data)
	b.PutUint16(buf[2+len(i.Data):], itemEnd)
}


func DecodeMsgPayload(buf []byte) (*MsgPayload, error) {
	var p MsgPayload
	//err := proto.Unmarshal(buf, &p)
	//if err != nil {
	//	log.Errorf("err:%v", err)
	//	return nil, err
	//}
	return &p, nil
}

type MsgPayload struct {
	//Ctx                  *MsgPayload_Context `protobuf:"bytes,1,opt,name=ctx" json:"ctx,omitempty"`
	MsgId string `protobuf:"bytes,2,opt,name=msg_id,json=msgId" json:"msg_id,omitempty"`
	Data  []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`

	DelayType uint32
	DelayValue uint32
}
