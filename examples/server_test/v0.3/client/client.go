package main

import (
	"fmt"
	"gmr/gmr-server/inet"
	"io"
	"net"
	"time"
)

func main() {
	fmt.Println("client start")
	time.Sleep(time.Second)
	//连接服务器得到conn连接
	conn,err := net.Dial("tcp","127.0.0.1:8089")
	if err != nil {
		fmt.Println("client start error",err)
		return
	}
	//连接调用write，写数据
	for{
		//发送封包消息
		dp := inet.NewDataPack()
		binaryMsg ,err := dp.Pack(inet.NewMsgPackage(1,[]byte("gmr-server v0.2 client")))
		if err != nil {
			fmt.Println("pack err",err)
			return
		}
		if _,err := conn.Write(binaryMsg);err != nil {
			fmt.Println("write err",err)
			return
		}

		//第一次读头
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("read head err", err)
			break
		}
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err", err)
			break
		}
		if msgHead.GetMsgLen() > 0 {
			//msg中有数据，第二次读内容
			msg := msgHead.(*inet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack err", err)
				break
			}
			fmt.Println("msgID", msg.Id, "msgLen", msg.DataLen, "msgData", string(msg.Data))
		}

		time.Sleep(time.Second)
	}
}