package iface

//将客户端连接信息，请求的数据包装到了Requst
type Request interface {
	//获取请求的连接信息
	GetConnection() Connection
	//获取请求的消息的数据
	GetData() []byte
	//获取请求的消息的ID
	GetMsgID() uint32
}
