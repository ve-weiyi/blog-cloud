package remark

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/messagerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRemarkReviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新留言
func NewUpdateRemarkReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRemarkReviewLogic {
	return &UpdateRemarkReviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRemarkReviewLogic) UpdateRemarkReview(req *types.RemarkReviewReq) (resp *types.BatchResp, err error) {
	in := &messagerpc.UpdateRemarkReviewReq{
		Ids:      req.Ids,
		IsReview: req.IsReview,
	}

	out, err := l.svcCtx.MessageRpc.UpdateRemarkReview(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
