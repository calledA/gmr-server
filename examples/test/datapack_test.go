package test

import (
	"fmt"
	"gmr/gmr-server/inet"
	"io"
	"net"
	"testing"
	"time"
)

func TestDataPack(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("listener server err", err)
	}

	//模拟服务端
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("accept err", err)
				continue
			}
			go func(conn net.Conn) {
				//处理客户端
				dp := inet.NewDataPack()
				for {
					//第一次读头
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
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
						//msg中有数据
						//第二次读内容
						msg := msgHead.(*inet.Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack err", err)
							break
						}
						fmt.Println("msgID", msg.Id, "msgLen", msg.DataLen, "msgData", msg.Data)
					}
				}
			}(conn)
		}
	}()

	//模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err", err)
	}
	//创建封包对象
	dp := inet.NewDataPack()

	//模拟一个msg1包
	msg1 := &inet.Message{
		Id:      1,
		DataLen: 3,
		Data:    []byte{'z', 'w', 'h'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("datapack err", err)
		return
	}
	//模拟一个msg2包
	msg2 := &inet.Message{
		Id:      2,
		DataLen: 4,
		Data:    []byte{'w', 'l', 'e', 'i'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("datapack err", err)
		return
	}
	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)

	//客户端阻塞
	// select {}
	time.Sleep(2 * time.Second)
}
