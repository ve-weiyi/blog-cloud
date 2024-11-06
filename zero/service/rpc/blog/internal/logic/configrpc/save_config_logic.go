package configrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/configrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveConfigLogic {
	return &SaveConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 保存配置
func (l *SaveConfigLogic) SaveConfig(in *configrpc.SaveConfigReq) (*configrpc.EmptyResp, error) {
	// 查找
	result, err := l.svcCtx.TWebsiteConfigModel.FindOneByKey(l.ctx, in.ConfigKey)
	if err != nil {
		return nil, err
	}

	// 修改
	result.Config = in.ConfigValue
	_, err = l.svcCtx.TWebsiteConfigModel.Save(l.ctx, result)
	if err != nil {
		return nil, err
	}

	return &configrpc.EmptyResp{}, nil
}
