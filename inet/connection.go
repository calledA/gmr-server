package inet

import (
	"context"
	"errors"
	"fmt"
	"gmr/gmr-server/iface"
	"io"
	"net"
)

//连接模块
type Connection struct {
	//TCP套接字
	Conn *net.TCPConn
	//连接ID
	ConnID uint32
	//当前连接状态
	isClosed bool
	//最大包字节数
	MaxPackageSize int
	//当前连接是否退出的channel
	ExitChan chan bool
	//绑定msgID处理对应API关系
	MsgHandler iface.IMsgHandler
}

//初始化
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler iface.IMsgHandler) iface.IConnection {
	return &Connection{
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
	}
}

//启动连接
func (c *Connection) Start() {
	fmt.Println("conn start connID:", c.ConnID)
	//启动当前连接读数据的业务
	go c.StartReader()
	//TODO 启动当前连接写数据的业务
}

func (c *Connection) StartReader() {
	fmt.Println("start reader is running")
	defer c.Stop()
	for {
		//读取客户端数据到buf，最大512
		// buf := make([]byte, config.GetPackageSize())
		// _, err := c.Conn.Read(buf)
		// if err != nil {
		// 	fmt.Println("recv buf err", err)
		// 	continue
		// }

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

		go c.MsgHandler.DoMsgHandler(&req)
	}
}

//停止连接，结束连接状态
func (c *Connection) Stop() {
	fmt.Println("conn stop connID:", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true
	c.Conn.Close()
	//回收资源
	close(c.ExitChan)
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
	//将数据发送给客户端
	if _, err := c.Conn.Write(msg); err != nil {
		fmt.Println("wirte msg err,msgID = ", msgID)
		return errors.New("conn write err")
	}
	return nil
}

//将Message数据发送给远程TCP客户端(有缓冲)
func (c *Connection) SendBuffMsg(msgID uint32, data []byte) error {

	return nil

}

//设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {}

//获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	return nil, nil
}

//移除连接属性
func (c *Connection) RemoveProperty(key string) {}
