package permissionrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateApiLogic {
	return &UpdateApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新接口
func (l *UpdateApiLogic) UpdateApi(in *permissionrpc.UpdateApiReq) (*permissionrpc.UpdateApiResp, error) {
	entity, err := l.svcCtx.TApiModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	entity.ParentId = in.ParentId
	entity.Path = in.Path
	entity.Name = in.Name
	entity.Method = in.Method
	entity.Traceable = in.Traceable
	entity.Status = in.Status

	_, err = l.svcCtx.TApiModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &permissionrpc.UpdateApiResp{
		Api: convertApiOut(entity),
	}, nil
}
