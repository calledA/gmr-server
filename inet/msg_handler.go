package inet

import (
	"fmt"
	"gmr/gmr-server/config"
	"gmr/gmr-server/iface"
	"strconv"
)

//消息管理实现
type MsgHandler struct {
	//存放每个msgID与对应的处理方法
	Apis map[uint32]iface.IRouter
	//worker取任务的消息队列
	TaskQueue []chan iface.IRequest
	//业务工作池的工作数量
	WorkerPoolSize uint32
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]iface.IRouter),
		TaskQueue:      make([]chan iface.IRequest, config.GetWorkerPoolSize()),
		WorkerPoolSize: config.GetWorkerPoolSize(), //从全局配置中获取
	}
}

//马上以非阻塞方式处理消息
func (mh *MsgHandler) DoMsgHandler(request iface.IRequest) {
	//从requst找到msgID，通过msgID处理对应router
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("apisID not found", request.GetMsgID())
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
func (mh *MsgHandler) StartWorkPool() {
	//根据workPoolSize分别开启worker
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//启动对应数量的worker，每个work的最大chan长度由GetMaxWorkerTaskLen()决定
		mh.TaskQueue[i] = make(chan iface.IRequest, config.GetMaxWorkerTaskLen())
		go mh.SendOneTask(i, mh.TaskQueue[i])
	}
}

//将消息交给TaskQueue，由work进行处理
func (mh *MsgHandler) SendMsgToTaskQueue(request iface.IRequest) {
	//将消息平均分配到worker队列，进行轮询
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize

	//将消息发送给对应的worker的TaskQueue
	mh.TaskQueue[workerID] <- request
}

func (mh *MsgHandler) SendOneTask(workID int, taskQueue chan iface.IRequest) {
	//不断阻塞等待对应消息队列的消息
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}
