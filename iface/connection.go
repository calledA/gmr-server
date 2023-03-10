package iface

import (
	"context"
	"net"
)

//连接接口
type IConnection interface {
	//启动连接
	Start()
	//停止连接，结束连接状态
	Stop()
	//返回ctx，用于用户自定义go程获取连接退出状态
	Context() context.Context
	//获取原始socket tcpConn
	GetTCPConnection() *net.TCPConn
	//获取连接Id
	GetConnID() uint32
	//获取远程客户端地址
	RemoteAddr() net.Addr
	//将Message数据发送给远程TCP客户端(无缓冲)
	SendMsg(msgID uint32, data []byte) error
}

type HandleFunc func(*net.TCPConn,[]byte,int) error