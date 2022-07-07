package call

type ConsumeReq struct {
	CreatedAt uint64  `protobuf:"varint,1,opt,name=created_at,json=createdAt" json:"created_at,omitempty"`
	RetryCnt  uint32 `protobuf:"varint,2,opt,name=retry_cnt,json=retryCnt" json:"retry_cnt,omitempty"`
	Data      []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	CorpId    uint32 `protobuf:"varint,4,opt,name=corp_id,json=corpId" json:"corp_id,omitempty"`
	AppId     uint32 `protobuf:"varint,5,opt,name=app_id,json=appId" json:"app_id,omitempty"`
	MsgId     string `protobuf:"bytes,6,opt,name=msg_id,json=msgId" json:"msg_id,omitempty"`
	// @desc: 服务端最大重试次数，业务可以根据这个来做最后失败逻辑
	MaxRetryCnt uint32 `protobuf:"varint,7,opt,name=max_retry_cnt,json=maxRetryCnt" json:"max_retry_cnt,omitempty"`
	Timeout     uint32
}
