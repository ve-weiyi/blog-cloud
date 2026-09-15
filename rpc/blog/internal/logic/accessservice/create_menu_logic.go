package accessservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/accessrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type CreateMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建菜单
func (l *CreateMenuLogic) CreateMenu(in *accessrpc.CreateMenuRequest) (*accessrpc.CreateMenuResponse, error) {
	menu := convertMenuIn(in)

	_, err := l.svcCtx.TMenuModel.Insert(l.ctx, menu)
	if err != nil {
		return nil, err
	}

	return &accessrpc.CreateMenuResponse{
		Id: menu.Id,
	}, nil
}
