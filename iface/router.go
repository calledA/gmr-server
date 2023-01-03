package iface

//使用框架者给该链接自定的 处理业务方法
type IRouter interface {
	//处理conn业务前的钩子方法
	PreHandle(request IRequest)
	//处理conn业务中的钩子方法
	Handle(request IRequest)
	//处理conn业务后的钩子方法
	PostHandle(request IRequest)
}
