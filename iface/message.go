package iface

//请求的消息封装到Message中，抽象接口层
type Message interface {
	//获取消息数据段长度
	GetDataLen() uint32
	//获取消息ID
	GetMsgID() uint32
	//获取消息内容
	GetData() []byte
	//设计消息ID
	SendMsgID(uint32)
	//设计消息内容
	SetData([]byte)
	//设置消息数据段长度
	SetDataLen(uint32)
}
