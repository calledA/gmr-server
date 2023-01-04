package iface

type IConnManager interface {
	//添加连接
	Add(conn IConnection)
	//删除连接
	Remove(conn IConnection)
	//利用ConnID获取连接
	Get(connID uint32) (IConnection, error)
	//获取当前连接数
	Len() int
	//删除并停止所有连接
	ClearConn()
}
