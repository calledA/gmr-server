package inet

import "gmr/gmr-server/iface"

type Request struct {
	//与客户端建立的连接
	conn iface.Connection
	//客户端请求的数据
	data []byte
}

//获取请求的连接信息
func (r *Request) GetConnection() iface.Connection {
	return r.conn
}
//获取请求的消息的数据
func (r *Request) GetData() []byte {
	return r.data
}
//获取请求的消息的ID
func (r *Request) GetMsgID() uint32 {
	return 0
}
