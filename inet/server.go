package inet

import (
	"fmt"
	"gmr/gmr-server/config"
	"gmr/gmr-server/iface"
	"net"
)

//Server接口的实现，定义一个server的服务模块
type Server struct {
	//服务器名称
	Name string
	//服务器IP版本
	IPVersion string
	//服务器的IP
	IP string
	//服务器端口
	Port int
	//服务器版本
	Version string
	//当前Server消息模块，绑定msgID处理对应API关系
	MsgHandler iface.IMsgHandler
	//server的连接管理器
	ConnManager iface.IConnManager
	//该server建立连接之前自动调用的hook
	OnConnStart func(conn iface.IConnection)
	//该server销毁连接之前自动调用的hook
	OnConnStop func(conn iface.IConnection)
}

//初始化Server方法
func NewServer() iface.IServer {
	conf := config.LoadServerConfig()
	return &Server{
		Name:        conf.Name,
		IPVersion:   conf.IPVersion,
		IP:          conf.IP,
		Port:        conf.Port,
		Version:     conf.Version,
		MsgHandler:  NewMsgHandler(),
		ConnManager: NewConnManager(),
	}
}

//启动服务器方法
func (s *Server) Start() {
	fmt.Printf("Server Listenner Addr:%s:%d \n", s.IP, s.Port)
	go func() {
		//开启消息队列及worker工作池
		s.MsgHandler.StartWorkPool()

		//获取TCP的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}
		//监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen tcp addr error:", err)
			return
		}
		var cid uint32 = 0

		fmt.Println("start server success")
		//阻塞等待连接，处理连接业务
		for {
			//阻塞，连接建立，阻塞停止
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept error:", err)
				continue
			}
			if s.ConnManager.Len() >= config.GetMaxConn() {
				//TODO 向客户端返回错误包
				conn.Close()
				continue
			}
			//将处理业务的方法与连接绑定
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			go dealConn.Start()
		}
	}()
}

//停止服务器方法
func (s *Server) Stop() {
	fmt.Println("stop server")
	//将服务器资源、状态和连接等进行停止和回收
	s.ConnManager.ClearConn()
}

//开启业务服务方法
func (s *Server) Serve() {
	//启动服务器业务
	s.Start()

	//TODO：服务器启动之后的额外业务

	//阻塞状态
	select {}
}

//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(msgID uint32, router iface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
}

//得到链接管理
func (s *Server) GetConnMgr() iface.IConnManager {
	return s.ConnManager
}

// 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(iface.IConnection)) {
	s.OnConnStart = hookFunc
}

//设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(iface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn iface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("call on CallOnConnStart")
		s.OnConnStart(conn)
	}
}

//调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn iface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("call on CallOnConnStop")
		s.OnConnStop(conn)
	}
}

// func (s *Server) Packet() iface.DataPack         {
// 	return iface.DataPack{}
// }
