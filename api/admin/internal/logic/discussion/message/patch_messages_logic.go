package message

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/discussionservice"
)

type PatchMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量更新留言状态
func NewPatchMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PatchMessagesLogic {
	return &PatchMessagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PatchMessagesLogic) PatchMessages(req *types.PatchMessagesReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.DiscussionService.BatchPatchMessages(l.ctx, &discussionservice.BatchPatchMessagesRequest{
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
