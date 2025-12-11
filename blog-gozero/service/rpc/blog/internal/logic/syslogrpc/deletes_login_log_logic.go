package syslogrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesLoginLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletesLoginLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesLoginLogLogic {
	return &DeletesLoginLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除登录记录
func (l *DeletesLoginLogLogic) DeletesLoginLog(in *syslogrpc.IdsReq) (*syslogrpc.BatchResp, error) {
	rows, err := l.svcCtx.TLoginLogModel.Deletes(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &syslogrpc.BatchResp{
		SuccessCount: rows,
	}, nil
}
