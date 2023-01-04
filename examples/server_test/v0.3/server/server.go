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
//创建链接之后执行钩子函数
func DoConnectionBegin(conn iface.IConnection) {
	fmt.Println("===> DoConnectionBegin is Called ... ")
	if err := conn.SendMsg(202, []byte("DoConnection BEGIN")); err != nil {
		fmt.Println(err)
	}
}

//链接断开之前的需要执行的函数
func DoConnectionLost(conn iface.IConnection) {
	fmt.Println("===> DoConnectionLost is Called ...")
	fmt.Println("conn ID = ", conn.GetConnID(), " is Lost...")
}

func main() {
	s := inet.NewServer()

	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)
	
	//添加自定义router
	s.AddRouter(0,&PingRouter{})
	s.AddRouter(1,&HelloRouter{})
	s.Serve()
}
