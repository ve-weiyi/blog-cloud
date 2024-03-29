package svc

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/repository"
	"github.com/ve-weiyi/ve-blog-golang/server/api/repository/svc"
	"github.com/ve-weiyi/ve-blog-golang/server/config"
	"github.com/ve-weiyi/ve-blog-golang/server/global"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/captcha"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/glog"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/jjwt"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/rabbitmq"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/rbac"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/upload"
)

// 注册需要用到的gorm、redis、model
type ServiceContext struct {
	*repository.AppRepository

	Config         *config.Config
	Log            *glog.Glogger
	Token          *jjwt.JwtToken
	RBAC           *rbac.CachedEnforcer
	Captcha        *captcha.CaptchaRepository
	EmailPublisher rabbitmq.MessagePublisher
	Uploader       upload.Uploader
}

func NewServiceContext(cfg *config.Config) *ServiceContext {
	ctx := svc.NewRepositoryContext(cfg)
	repo := repository.NewRepository(ctx)
	if repo == nil {
		panic("repository cannot be null")
	}

	return &ServiceContext{
		AppRepository:  repo,
		Config:         cfg,
		Log:            global.LOG,
		Token:          global.JWT,
		Captcha:        captcha.NewCaptchaRepository(),
		EmailPublisher: global.EmailMQ,
		Uploader:       global.Uploader,
	}
}
