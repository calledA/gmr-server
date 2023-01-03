package main

import (
	"fmt"
	"gmr/gmr-server/iface"
	"gmr/gmr-server/inet"
)

type PingRouter struct {
	inet.BaseRouter
}

//处理conn业务中的钩子方法
func (br *PingRouter) Handle(request iface.IRequest) {
	fmt.Println("call router handler")
	//先读取客户端数据，再回写ping
	fmt.Println("msgID:", request.GetMsgID(), " msgData", string(request.GetData()))
	err := request.GetConnection().SendMsg(200, []byte("ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloRouter struct {
	inet.BaseRouter
}

//处理conn业务中的钩子方法
func (br *HelloRouter) Handle(request iface.IRequest) {
	fmt.Println("call router handler")
	//先读取客户端数据，再回写ping
	fmt.Println("msgID:", request.GetMsgID(), " msgData", string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("hello world"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := inet.NewServer()
	//添加自定义router
	s.AddRouter(0,&PingRouter{})
	s.AddRouter(1,&HelloRouter{})
	s.Serve()
}
