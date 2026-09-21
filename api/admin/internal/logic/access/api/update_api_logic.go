package api

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/accessservice"
)

type UpdateApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新API
func NewUpdateApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateApiLogic {
	return &UpdateApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateApiLogic) UpdateApi(req *types.UpdateApiReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.AccessService.UpdateApi(l.ctx, &accessservice.UpdateApiRequest{
		Id:        req.Id,
		ParentId:  req.ParentId,
		Name:      req.Name,
		Path:      req.Path,
		Method:    req.Method,
		Traceable: req.Traceable,
		Status:    req.Status,
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
