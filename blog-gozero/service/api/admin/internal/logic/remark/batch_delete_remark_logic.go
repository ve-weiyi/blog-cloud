package remark

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/messagerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteRemarkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除留言
func NewBatchDeleteRemarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteRemarkLogic {
	return &BatchDeleteRemarkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteRemarkLogic) BatchDeleteRemark(req *types.IdsReq) (resp *types.BatchResp, err error) {
	in := &messagerpc.IdsReq{
		Ids: req.Ids,
	}

	out, err := l.svcCtx.MessageRpc.DeletesRemark(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
