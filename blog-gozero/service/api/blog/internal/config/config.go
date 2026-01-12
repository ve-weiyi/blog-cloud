package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/ve-weiyi/ve-blog-golang/pkg/kit/oss"
)

type Config struct {
	rest.RestConf
	RedisConf    RedisConf
	AiProxyConf  AiProxyConf
	UploadConfig *oss.Config

	BlogRpcConf zrpc.RpcClientConf
}

// redis缓存配置
type RedisConf struct {
	Host     string `json:"host,default=127.0.0.1"`
	Port     string `json:"port,default=6379"`
	Password string `json:"password,default=''"`
	DB       int    `json:"db,default=0"`
}

// ai代理配置
type AiProxyConf struct {
	ApiHost string `json:"api_host,default=https://api.openai.com"`
	ApiKey  string `json:"api_key,default=''"`
	Model   string `json:"model,default=gpt-3.5-turbo"`
}
