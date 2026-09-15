package menu

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/accessservice"
)

type CleanMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 清空菜单列表
func NewCleanMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CleanMenuLogic {
	return &CleanMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CleanMenuLogic) CleanMenu(req *types.CleanMenuReq) (resp *types.CleanMenuResp, err error) {
	out, err := l.svcCtx.AccessService.CleanMenus(l.ctx, &accessservice.CleanMenusRequest{})
	if err != nil {
		return nil, err
	}

	return &types.CleanMenuResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
