package message

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/discussionservice"
)

type UpdateMessageStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量更新留言状态
func NewUpdateMessageStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMessageStatusLogic {
	return &UpdateMessageStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMessageStatusLogic) UpdateMessageStatus(req *types.UpdateMessageStatusReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.DiscussionService.BatchUpdateMessageStatus(l.ctx, &discussionservice.BatchUpdateMessageStatusRequest{
		Ids:    req.Ids,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
