package configrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/configrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

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
func (l *SaveConfigLogic) SaveConfig(in *configrpc.SaveConfigReq) (*configrpc.SaveConfigResp, error) {
	// 修改
	entity := &model.TWebsiteConfig{
		Id:     0,
		Key:    in.ConfigKey,
		Config: in.ConfigValue,
	}

	// 查找
	result, err := l.svcCtx.TWebsiteConfigModel.FindOneByKey(l.ctx, in.ConfigKey)
	if result != nil {
		entity.Id = result.Id
	}

	_, err = l.svcCtx.TWebsiteConfigModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &configrpc.SaveConfigResp{}, nil
}
