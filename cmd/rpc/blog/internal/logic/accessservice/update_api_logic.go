package accessservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/accessrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/model"
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

	// 更新后回读：返回体需要包含未在请求中的字段（创建时间等）
	entity, err := l.svcCtx.TApiModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &accessrpc.UpdateApiResponse{Success: true, Api: convertApiOut(entity)}, nil
}
