package permissionrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddApiLogic {
	return &AddApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建接口
func (l *AddApiLogic) AddApi(in *permissionrpc.AddApiReq) (*permissionrpc.AddApiResp, error) {
	entity := convertApiIn(in)
	_, err := l.svcCtx.TApiModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &permissionrpc.AddApiResp{
		Api: convertApiOut(entity),
	}, nil
}

func convertApiIn(in *permissionrpc.AddApiReq) (out *model.TApi) {
	out = &model.TApi{
		Id:        in.Id,
		ParentId:  in.ParentId,
		Name:      in.Name,
		Path:      in.Path,
		Method:    in.Method,
		Traceable: in.Traceable,
		Status:    in.Status,
	}

	return out
}

func convertApiOut(in *model.TApi) (out *permissionrpc.Api) {
	out = &permissionrpc.Api{
		Id:        in.Id,
		ParentId:  in.ParentId,
		Name:      in.Name,
		Path:      in.Path,
		Method:    in.Method,
		Traceable: in.Traceable,
		Status:    in.Status,
		CreatedAt: in.CreatedAt.UnixMilli(),
		UpdatedAt: in.UpdatedAt.UnixMilli(),
	}

	return out
}
