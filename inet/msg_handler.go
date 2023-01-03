package inet

import (
	"fmt"
	"gmr/gmr-server/iface"
	"strconv"
)

//消息管理实现
type MsgHandler struct {
	Apis map[uint32]iface.IRouter
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]iface.IRouter),
	}
}

//马上以非阻塞方式处理消息
func (mh *MsgHandler) DoMsgHandler(request iface.IRequest) {
	//从requst找到msgID，通过msgID处理对应router
	handler,ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("apisID not found",request.GetMsgID())
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//为消息添加具体处理逻辑
func (mh *MsgHandler) AddRouter(msgID uint32, router iface.IRouter) {
	//判断当前Apis的处理方法是否存在
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat api,msgID = " + strconv.Itoa(int(msgID)))
	}
	mh.Apis[msgID] = router
}

//启动worker工作池
func (mh *MsgHandler) StartWorkPool() {}

//将消息交给TaskQueue，由work进行处理
func (mh *MsgHandler) SendMsgToTaskQueue(request iface.IRequest) {}
