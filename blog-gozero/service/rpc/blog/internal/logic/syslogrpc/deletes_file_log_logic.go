package syslogrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesFileLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletesFileLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesFileLogLogic {
	return &DeletesFileLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除文件记录
func (l *DeletesFileLogLogic) DeletesFileLog(in *syslogrpc.DeletesFileLogReq) (*syslogrpc.DeletesFileLogResp, error) {
	rows, err := l.svcCtx.TFileLogModel.Deletes(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &syslogrpc.DeletesFileLogResp{
		SuccessCount: rows,
	}, nil
}
