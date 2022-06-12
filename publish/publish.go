package publish


type PubReq struct {
	// @v: required
	TopicName string `protobuf:"bytes,1,opt,name=topic_name,json=topicName" json:"topic_name,omitempty" validate:"required"`
	// @v: required
	Data  []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty" validate:"required"`
	ReqId string `protobuf:"bytes,5,opt,name=req_id,json=reqId" json:"req_id,omitempty"`
	// @desc: 可选的 msg_id
	MsgId  string `protobuf:"bytes,10,opt,name=msg_id,json=msgId" json:"msg_id,omitempty"`
	Hash   uint32 `protobuf:"varint,13,opt,name=hash" json:"hash,omitempty"`
	Group  string `protobuf:"bytes,20,opt,name=group" json:"group,omitempty"`
	Node   string `protobuf:"bytes,21,opt,name=node" json:"node,omitempty"`
	// @desc: delay 仅对 barrier mode 生效
	//  delay_type == DelayTypeRelate 时，delay_value 的单位是 ms
	DelayType  uint32 `protobuf:"varint,30,opt,name=delay_type,json=delayType" json:"delay_type,omitempty"`
	DelayValue uint32 `protobuf:"varint,31,opt,name=delay_value,json=delayValue" json:"delay_value,omitempty"`
	// @desc: [仅串行化模式有效] 范围是 0-3，数值越大，越优先执行
	Priority             uint32   `protobuf:"varint,32,opt,name=priority" json:"priority,omitempty"`
}
