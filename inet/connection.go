package inet

import (
	"context"
	"errors"
	"fmt"
	"gmr/gmr-server/config"
	"gmr/gmr-server/iface"
	"io"
	"net"
)

//连接模块
type Connection struct {
	//属于哪个server
	TCPServer iface.IServer
	//TCP套接字
	Conn *net.TCPConn
	//连接ID
	ConnID uint32
	//当前连接状态
	isClosed bool
	//最大连接数
	MaxConn int
	//最大包字节数
	MaxPackageSize int
	//当前连接是否退出的channel,由reader告知writer退出
	ExitChan chan bool
	//无缓冲管道，用于读写的通道消息通信
	msgChan chan []byte
	//绑定msgID处理对应API关系
	MsgHandler iface.IMsgHandler
}

//初始化
func NewConnection(server iface.IServer, conn *net.TCPConn, connID uint32, msgHandler iface.IMsgHandler) iface.IConnection {
	c := &Connection{
		TCPServer:  server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
	}
	c.TCPServer.GetConnMgr().Add(c)
	return c
}

//启动连接
func (c *Connection) Start() {
	fmt.Println("conn start connID:", c.ConnID)
	//启动当前连接读数据的业务
	go c.StartReader()
	//启动当前连接写数据的业务
	go c.StartWriter()

	//使用自定义的hook函数
	c.TCPServer.CallOnConnStart(c)
}

//读消息goroutine，专门接收客户端消息
func (c *Connection) StartReader() {
	defer fmt.Println("[reader is stop]")
	fmt.Println("[start reader is running]")
	defer c.Stop()
	for {
		//创建拆包解包对象
		dp := NewDataPack()
		//读取客户端的msg header
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			break
		}
		//拆包，得到msgID和msgDataLen放在消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		//根据dataLen，再次读取data放入msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg error", err)
				break
			}
		}
		msg.SetData(data)

		//得到当前conn数据的request的请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if config.GetWorkerPoolSize() > 0 {
			//已经开启工作池，通过工作池处理业务
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//没有开启工作池，每个连接通过一个goroutine处理业务
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

//写消息goroutine，专门向客户端发送消息
func (c *Connection) StartWriter() {
	defer fmt.Println("[writer is stop]")
	fmt.Println("[writer is running]")
	//不断阻塞等待channel消息，进行写消息
	for {
		select {
		case data := <-c.msgChan:
			//如果有数据写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data error", err)
				return
			}
		case <-c.ExitChan:
			//代表reader退出，此时writer也要退出
			return
		}
	}
}

//停止连接，结束连接状态
func (c *Connection) Stop() {
	fmt.Println("conn stop connID:", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true
	//调用自定义的hook
	c.TCPServer.CallOnConnStop(c)
	//关闭socket连接
	c.Conn.Close()
	//告知writer，关闭连接
	c.ExitChan <- true
	//将当前连接从connManager中摘除
	c.TCPServer.GetConnMgr().Remove(c)
	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

//返回ctx，用于用户自定义go程获取连接退出状态
func (c *Connection) Context() context.Context {
	return context.Background()
}

//获取原始socket tcpConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取连接Id
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端地址
func (c *Connection) RemoteAddr() net.Addr {
	return nil
}

//将Message数据发送给远程TCP客户端(无缓冲)
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when send msg")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("pack msg err,msgID = ", msgID)
		return errors.New("pack msg err")
	}
	//将数据发送给msgChan
	c.msgChan <- msg
	return nil
}


