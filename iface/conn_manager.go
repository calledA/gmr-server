package iface

type ConnManager interface {
	//添加连接
	Add(conn Connection)
	//删除连接
	Remove(conn Connection)
	//利用ConnID获取连接
	Get(connID uint32) (Connection, error)
	//获取当前连接
	Len() int
	//删除并停止所有连接
	ClearConn()
}
