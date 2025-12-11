package svc

import (
	"encoding/json"
	"fmt"

	"github.com/go-openapi/loads"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/permissionx"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/tokenx"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/docs"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/common/stomphook"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/config"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/middleware"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/configrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/resourcerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/talkrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/websiterpc"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/oss"
	"github.com/ve-weiyi/ve-blog-golang/stompws/logws"
	"github.com/ve-weiyi/ve-blog-golang/stompws/server/client"
)

type ServiceContext struct {
	Config config.Config

	AccountRpc    accountrpc.AccountRpc
	PermissionRpc permissionrpc.PermissionRpc
	ArticleRpc    articlerpc.ArticleRpc
	MessageRpc    messagerpc.MessageRpc
	ResourceRpc   resourcerpc.ResourceRpc
	TalkRpc       talkrpc.TalkRpc
	WebsiteRpc    websiterpc.WebsiteRpc
	ConfigRpc     configrpc.ConfigRpc
	SyslogRpc     syslogrpc.SyslogRpc

	Redis            *redis.Redis
	Uploader         oss.OSS
	TokenHolder      tokenx.TokenHolder
	PermissionHolder permissionx.PermissionHolder

	StompHubServer *client.StompHubServer

	AdminToken   rest.Middleware
	Permission   rest.Middleware
	OperationLog rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	var options []zrpc.ClientOption
	options = append(options)

	rds, err := ConnectRedis(c.RedisConf)
	if err != nil {
		panic(err)
	}

	uploader := oss.NewQiniu(c.UploadConfig)

	th := tokenx.NewJwtTokenHolder(c.Name, c.Name, rds)

	doc, err := loads.Analyzed(json.RawMessage(docs.Docs), "")
	if err != nil {
		panic(err)
	}

	accountRpc := accountrpc.NewAccountRpc(zrpc.MustNewClient(c.BlogRpcConf, options...))
	permissionRpc := permissionrpc.NewPermissionRpc(zrpc.MustNewClient(c.BlogRpcConf, options...))
	articleRpc := articlerpc.NewArticleRpc(zrpc.MustNewClient(c.BlogRpcConf, options...))
	messageRpc := messagerpc.NewMessageRpc(zrpc.MustNewClient(c.BlogRpcConf, options...))
	resourceRpc := resourcerpc.NewResourceRpc(zrpc.MustNewClient(c.BlogRpcConf, options...))
	talkRpc := talkrpc.NewTalkRpc(zrpc.MustNewClient(c.BlogRpcConf, options...))
	websiteRpc := websiterpc.NewWebsiteRpc(zrpc.MustNewClient(c.BlogRpcConf, options...))
	configRpc := configrpc.NewConfigRpc(zrpc.MustNewClient(c.BlogRpcConf, options...))
	syslogRpc := syslogrpc.NewSyslogRpc(zrpc.MustNewClient(c.BlogRpcConf, options...))

	ph := permissionx.NewMemoryHolder(permissionRpc)
	err = ph.LoadPolicy()
	if err != nil {
		logx.Infof("load permission policy fail: %v", err)
	}

	hub := client.NewStompHubServer(
		client.WithEventHooks(
			stomphook.NewChatRoomEventHook(),
			stomphook.NewOnlineEventHook(),
		),
		client.WithAuthenticator(stomphook.NewJwtAuthenticator(th)),
		client.WithLogger(logws.NewDefaultLogger()),
	)

	return &ServiceContext{
		Config:        c,
		AccountRpc:    accountRpc,
		PermissionRpc: permissionRpc,
		ArticleRpc:    articleRpc,
		MessageRpc:    messageRpc,
		ResourceRpc:   resourceRpc,
		TalkRpc:       talkRpc,
		WebsiteRpc:    websiteRpc,
		ConfigRpc:     configRpc,
		SyslogRpc:     syslogRpc,

		Redis:            rds,
		Uploader:         uploader,
		TokenHolder:      th,
		PermissionHolder: ph,
		StompHubServer:   hub,

		AdminToken: middleware.NewAdminTokenMiddleware(th).Handle,
		Permission: middleware.NewSimpleMiddleware().Handle, // 不使用接口权限控制
		//Permission:       middleware.NewPermissionMiddleware(ph).Handle,
		OperationLog: middleware.NewOperationLogMiddleware(doc.Spec(), syslogRpc, permissionRpc).Handle,
	}
}

func ConnectRedis(c config.RedisConf) (*redis.Redis, error) {
	address := c.Host + ":" + c.Port
	client, err := redis.NewRedis(redis.RedisConf{
		Host: address,
		Type: redis.NodeType,
		Pass: c.Password,
		Tls:  false,
	})

	if err != nil {
		return nil, fmt.Errorf("redis 连接失败: %v", err)
	}

	return client, nil
}
