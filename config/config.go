package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	//Server信息
	Server *Server
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
	//最大连接数
	MaxConn int
	//最大包字节
	MaxPackageSize int
}

var config *Config

func LoadServerConfig() Server {
	return *config.Server
}

func GetPackageSize() int {
	return config.Server.MaxPackageSize
}

func (config *Config) LoadConfig() {
	conf := viper.New()
	conf.AddConfigPath("../../config") //设置读取的文件路径
	conf.SetConfigName("app")          //设置读取的文件名
	conf.SetConfigType("yml")          //设置文件的类型
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
	config.Server.MaxConn = conf.GetInt("server.max_conn")
	config.Server.MaxPackageSize = conf.GetInt("server.max_package_size")
}

func init() {
	config = &Config{
		Server: &Server{
			Name:           "GMR-SERVER",
			Version:        "V0.4",
			Port:           8089,
			IP:             "0.0.0.0",
			MaxConn:        1000,
			MaxPackageSize: 4096,
		},
	}
	config.LoadConfig()
}
