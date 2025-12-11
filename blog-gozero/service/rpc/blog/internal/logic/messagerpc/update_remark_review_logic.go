package messagerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRemarkReviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRemarkReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRemarkReviewLogic {
	return &UpdateRemarkReviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新留言审核状态
func (l *UpdateRemarkReviewLogic) UpdateRemarkReview(in *messagerpc.UpdateRemarkReviewReq) (*messagerpc.BatchResp, error) {
	rows, err := l.svcCtx.TRemarkModel.Updates(l.ctx, map[string]interface{}{
		"is_review": in.IsReview,
	}, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &messagerpc.BatchResp{
		SuccessCount: rows,
	}, nil
}
