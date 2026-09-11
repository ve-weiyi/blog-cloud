package permissionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type CreateApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateApiLogic {
	return &CreateApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建 API
func (l *CreateApiLogic) CreateApi(in *permissionrpc.CreateApiRequest) (*permissionrpc.CreateApiResponse, error) {
	api := convertApiIn(in)

	_, err := l.svcCtx.TApiModel.Insert(l.ctx, api)
	if err != nil {
		return nil, err
	}

	return &permissionrpc.CreateApiResponse{
		Id: api.Id,
	}, nil
}
