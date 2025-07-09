package syslogrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesUploadLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletesUploadLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesUploadLogLogic {
	return &DeletesUploadLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除上传记录
func (l *DeletesUploadLogLogic) DeletesUploadLog(in *syslogrpc.IdsReq) (*syslogrpc.BatchResp, error) {
	rows, err := l.svcCtx.TUploadLogModel.Deletes(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &syslogrpc.BatchResp{
		SuccessCount: rows,
	}, nil
}
