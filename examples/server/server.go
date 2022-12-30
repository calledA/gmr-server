package main

import (
	"fmt"
	"gmr/gmr-server/iface"
	"gmr/gmr-server/inet"
)

type PingRouter struct {
	inet.BaseRouter
}

//处理conn业务前的钩子方法
func (br *PingRouter) PreHandler(request iface.Request) {
	fmt.Println("call router prehandler")
	_,err := request.GetConnection().GetTCPConnection().Write([]byte("prehandler ping"))
	if err != nil {
		fmt.Println("prehandler call back err",err)
	}
}

//处理conn业务中的钩子方法
func (br *PingRouter) Handle(request iface.Request) {
	fmt.Println("call router handler")
	_,err := request.GetConnection().GetTCPConnection().Write([]byte("handler ping"))
	if err != nil {
		fmt.Println("handler call back err",err)
	}
}

//处理conn业务后的钩子方法
func (br *PingRouter) PostHandle(request iface.Request) {
	fmt.Println("call router posthandler")
	_,err := request.GetConnection().GetTCPConnection().Write([]byte("posthandler ping"))
	if err != nil {
		fmt.Println("posthandler call back err",err)
	}
}

func main() {
	s := inet.NewServer()
	//添加自定义router
	s.AddRouter(&PingRouter{})
	s.Serve()
}