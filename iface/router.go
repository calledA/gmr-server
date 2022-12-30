package iface

//使用框架者给该链接自定的 处理业务方法
type Router interface {
	//处理conn业务前的钩子方法
	PreHandler(request Request)
	//处理conn业务中的钩子方法
	Handle(request Request)
	//处理conn业务后的钩子方法
	PostHandle(request Request)
}
