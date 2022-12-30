package inet

import "gmr/gmr-server/iface"

//实现router，先嵌入这个基类，然后对基类进行重写
type BaseRouter struct{}

/*
BaseRouter将Router接口空实现，自定义Router只需要继承BaseRouter，可以重写这几个方法，继承BaseRouter可以不用实现所有方法
*/
//处理conn业务前的钩子方法
func (br *BaseRouter) PreHandler(request iface.Request) {}

//处理conn业务中的钩子方法
func (br *BaseRouter) Handle(request iface.Request) {}

//处理conn业务后的钩子方法
func (br *BaseRouter) PostHandle(request iface.Request) {}
