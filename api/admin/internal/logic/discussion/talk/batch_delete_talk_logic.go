package talk

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/discussionservice"
)

type BatchDeleteTalkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除说说
func NewBatchDeleteTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteTalkLogic {
	return &BatchDeleteTalkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteTalkLogic) BatchDeleteTalk(req *types.DeleteTalkReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.DiscussionService.BatchDeleteTalks(l.ctx, &discussionservice.BatchDeleteTalksRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
