package message

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/discussionservice"
)

type BatchDeleteMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除留言
func NewBatchDeleteMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteMessageLogic {
	return &BatchDeleteMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteMessageLogic) BatchDeleteMessage(req *types.DeleteMessageReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.DiscussionService.BatchDeleteMessages(l.ctx, &discussionservice.BatchDeleteMessagesRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
