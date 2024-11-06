package configrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/configrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindConfigLogic {
	return &FindConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询配置
func (l *FindConfigLogic) FindConfig(in *configrpc.FindConfigReq) (*configrpc.FindConfigResp, error) {
	entity, err := l.svcCtx.TWebsiteConfigModel.FindOneByKey(l.ctx, in.ConfigKey)
	if err != nil {
		return nil, err
	}

	return &configrpc.FindConfigResp{
		ConfigValue: entity.Config,
	}, nil
}
