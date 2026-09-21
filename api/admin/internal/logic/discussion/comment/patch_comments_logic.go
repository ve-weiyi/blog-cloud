package comment

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/discussionservice"
)

type PatchCommentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量更新评论状态
func NewPatchCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PatchCommentsLogic {
	return &PatchCommentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PatchCommentsLogic) PatchComments(req *types.PatchCommentsReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.DiscussionService.BatchPatchComments(l.ctx, &discussionservice.BatchPatchCommentsRequest{
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
