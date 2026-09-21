package accessservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/accessrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
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

// 更新 API
func (l *UpdateApiLogic) UpdateApi(in *accessrpc.UpdateApiRequest) (*accessrpc.UpdateApiResponse, error) {
	data := &model.TApi{
		Id:        in.Id,
		ParentId:  in.ParentId,
		Path:      in.Path,
		Name:      in.Name,
		Method:    in.Method,
		Traceable: in.Traceable,
		Status:    in.Status,
	}
	_, err := l.svcCtx.TApiModel.Update(l.ctx, data)
	if err != nil {
		return nil, err
	}

	return &accessrpc.UpdateApiResponse{Success: true}, nil
}
