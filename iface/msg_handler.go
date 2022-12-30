package iface

//消息管理抽象层
type MsgHandler interface {
	//马上以非阻塞方式处理消息
	DoMsgHandler(request Request)
	//为消息添加具体处理逻辑
	AddRouter(msgID uint32,router Router)
	//启动worker工作池
	StartWorkPool()
	//将消息交给TaskQueue，由work进行处理
	SendMsgToTaskQueue(request Request)
}
