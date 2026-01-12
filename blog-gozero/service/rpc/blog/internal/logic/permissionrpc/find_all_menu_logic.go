package permissionrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindAllMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindAllMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindAllMenuLogic {
	return &FindAllMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查找所有菜单
func (l *FindAllMenuLogic) FindAllMenu(in *permissionrpc.FindAllMenuReq) (*permissionrpc.FindAllMenuResp, error) {
	result, err := l.svcCtx.TMenuModel.FindALL(l.ctx, "")
	if err != nil {
		return nil, err
	}

	out := &permissionrpc.FindAllMenuResp{}
	for _, item := range result {
		out.List = append(out.List, convertMenuOut(item))
	}

	return out, nil
}
