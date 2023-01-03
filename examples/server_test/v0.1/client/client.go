package main

import (
	"fmt"
	"gmr/gmr-server/config"
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
		_,err := conn.Write([]byte("hello world"))
		if err != nil {
			fmt.Println("write error",err)
			return
		}
		buf := make([]byte,config.GetPackageSize())
		cnt,err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error",err)
			return
		}
		fmt.Printf("call back:%s cnt:%d \n",buf,cnt)
		time.Sleep(time.Second)
	}
}