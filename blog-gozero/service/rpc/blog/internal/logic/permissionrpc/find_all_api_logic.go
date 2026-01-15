package permissionrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindAllApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindAllApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindAllApiLogic {
	return &FindAllApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查找所有接口
func (l *FindAllApiLogic) FindAllApi(in *permissionrpc.FindAllApiReq) (*permissionrpc.FindAllApiResp, error) {
	result, err := l.svcCtx.TApiModel.FindALL(l.ctx, "")
	if err != nil {
		return nil, err
	}

	out := &permissionrpc.FindAllApiResp{}
	for _, item := range result {
		out.List = append(out.List, convertApiOut(item))
	}

	return out, nil
}
