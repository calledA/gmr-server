package inet

import (
	"errors"
	"fmt"
	"gmr/gmr-server/iface"
	"sync"
)

type ConnManager struct {
	//管理连接的集合
	connections map[uint32]iface.IConnection
	//读写锁
	lock sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]iface.IConnection),
	}
}

//添加连接
func (cm *ConnManager) Add(conn iface.IConnection) {
	//共享资源加写锁
	cm.lock.Lock()
	defer cm.lock.Unlock()
	fmt.Println("ConnManager add", conn.GetConnID())

	//将conn添加到ConnManager中
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("Add()", len(cm.connections))
}

//删除连接
func (cm *ConnManager) Remove(conn iface.IConnection) {
	//共享资源加写锁
	cm.lock.Lock()
	defer cm.lock.Unlock()

	//删除map中的信息
	delete(cm.connections, conn.GetConnID())
	fmt.Println("remove connManager success", cm.Len())
}

//利用ConnID获取连接
func (cm *ConnManager) Get(connID uint32) (iface.IConnection, error) {
	//共享资源加读锁
	cm.lock.RLock()
	defer cm.lock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("could not found connID")
	}
}

//获取当前连接数
func (cm *ConnManager) Len() int {
	cm.lock.RLock()
	defer cm.lock.RUnlock()
	return len(cm.connections)
}

//删除并停止所有连接
func (cm *ConnManager) ClearConn() {
	//共享资源加写锁
	cm.lock.Lock()
	defer cm.lock.Unlock()

	//删除conn，并停止工作
	for id, conn := range cm.connections {
		//停止连接
		conn.Stop()
		//删除map
		delete(cm.connections, id)
	}

	fmt.Println("clear all connections success")
}
