package iface

const (
	//标准封包拆包方法
	StandardDataPack string = "standard_pack"
	//自定义封包拆包方法
	// ...

	//标准报文格式协议
	StandardMessage string = "standard_message"
)

/*
	封包数据和拆包数据
	直接面向TCP连接中的数据流,为传输数据添加头部信息，用于处理TCP粘包问题。
*/
type IDataPack interface {
	//获取包头长度
	GetHeadLen() uint32
	//封包方法，先写长度、类型，再写消息内容
	Pack(msg IMessage) ([]byte, error)
	//拆包方法,先读到包长度和包类型，再读取消息内容
	Unpack([]byte) (IMessage, error)
}

