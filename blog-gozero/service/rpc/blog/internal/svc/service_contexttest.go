package svc

import (
	"log"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/config"
)

func NewTestConfig() config.Config {
	return config.Config{
		RpcServerConf: zrpc.RpcServerConf{
			ServiceConf: service.ServiceConf{
				Mode: "dev",
				Log: logx.LogConf{
					Mode:     "console",
					Encoding: "plain",
					Path:     "logs",
				},
			},
		},
		MysqlConf: config.MysqlConf{
			Host:     "127.0.0.1",
			Port:     "3306",
			Username: "root",
			Password: "mysql7914",
			Dbname:   "blog-veweiyi",
			Config:   "charset=utf8mb4&parseTime=True&loc=Local",
		},
		RedisConf: config.RedisConf{
			DB:       0,
			Host:     "127.0.0.1",
			Port:     "6379",
			Password: "redis7914",
		},
		RabbitMQConf: config.RabbitMQConf{
			Host:     "127.0.0.1",
			Port:     "5672",
			Username: "veweiyi",
			Password: "rabbitmq7914",
		},
		EmailConf: config.EmailConf{
			Host:     "smtp.qq.com",
			Port:     465,
			Username: "",
			Password: "",
			Nickname: "",
			Deliver:  nil,
		},
		OauthConfList: nil,
	}
}

func NewTestServiceContext() *ServiceContext {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	//c := NewTestConfig()

	c := config.Config{}
	conf.MustLoad("../../etc/blog-rpc.yaml", &c)
	return NewServiceContext(c)
}
