package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	//Server信息
	Server *Server
	//Connection信息
	Connection *Connection
	//MsgHandler信息
	MsgHandler *MsgHandler
}

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
}

type Connection struct {
	//最大连接数
	MaxConn int
	//最大包字节
	MaxPackageSize uint32
}

type MsgHandler struct {
	//工作池大小
	WorkerPoolSize uint32
	//最大Task长度
	MaxWorkerTaskLen uint32
}

var config *Config

func LoadServerConfig() Server {
	return *config.Server
}

func GetPackageSize() uint32 {
	return config.Connection.MaxPackageSize
}

func GetMaxConn() int {
	return config.Connection.MaxConn
}

func GetWorkerPoolSize() uint32 {
	return config.MsgHandler.WorkerPoolSize
}

func GetMaxWorkerTaskLen() uint32 {
	return config.MsgHandler.MaxWorkerTaskLen
}

func (config *Config) LoadConfig() {
	conf := viper.New()
	conf.AddConfigPath("../../../../config") //设置读取的文件路径
	conf.SetConfigName("app")                //设置读取的文件名
	conf.SetConfigType("yml")                //设置文件的类型
	//尝试进行配置读取
	if err := conf.ReadInConfig(); err != nil {
		fmt.Println("read config error", err)
	}

	//获取最新配置
	config.Server.IP = conf.GetString("server.host")
	config.Server.Port = conf.GetInt("server.port")
	config.Server.IPVersion = conf.GetString("server.ip_version")
	config.Server.Name = conf.GetString("server.name")
	config.Server.Version = conf.GetString("server.version")
	config.Connection.MaxConn = conf.GetInt("server.max_conn")
	config.Connection.MaxPackageSize = conf.GetUint32("server.max_package_size")
	config.MsgHandler.WorkerPoolSize = conf.GetUint32("server.worker_pool_size")
	config.MsgHandler.MaxWorkerTaskLen = conf.GetUint32("server.max_worker_task_len")
}

func init() {
	config = &Config{
		Server: &Server{
			Name:           "GMR-SERVER",
			Version:        "V0.4",
			Port:           8089,
			IP:             "0.0.0.0",
		},
		Connection: &Connection{
			MaxConn:        1000,
			MaxPackageSize: 4096,
		},
		MsgHandler: &MsgHandler{
			WorkerPoolSize: 10,
			MaxWorkerTaskLen: 1024,
		},
	}
	config.LoadConfig()
}
