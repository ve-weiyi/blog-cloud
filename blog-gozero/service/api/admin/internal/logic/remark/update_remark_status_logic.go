package remark

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/messagerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRemarkStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新留言
func NewUpdateRemarkStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRemarkStatusLogic {
	return &UpdateRemarkStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRemarkStatusLogic) UpdateRemarkStatus(req *types.UpdateRemarkStatusReq) (resp *types.BatchResp, err error) {
	in := &messagerpc.UpdateRemarkStatusReq{
		Ids:    req.Ids,
		Status: req.Status,
	}

	out, err := l.svcCtx.MessageRpc.UpdateRemarkStatus(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
