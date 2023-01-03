package inet

import "gmr/gmr-server/iface"

type Request struct {
	//与客户端建立的连接
	conn iface.IConnection
	//客户端请求的数据
	msg iface.IMessage
}

//获取请求的连接信息
func (r *Request) GetConnection() iface.IConnection {
	return r.conn
}
//获取请求的消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}
//获取请求的消息的ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}
