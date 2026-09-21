package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/ve-weiyi/vkit/adapter/storagex"
)

type Config struct {
	rest.RestConf
	RedisConf     RedisConf
	StorageConfig storagex.StorageConfig

	AppRpcConf zrpc.RpcClientConf
}

// redis缓存配置
type RedisConf struct {
	Host     string `json:"host,default=127.0.0.1"`
	Port     string `json:"port,default=6379"`
	Password string `json:"password,default=''"`
	DB       int    `json:"db,default=0"`
}
