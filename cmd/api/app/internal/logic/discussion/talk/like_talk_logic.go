package talk

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/discussionservice"
)

type LikeTalkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 点赞说说
func NewLikeTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeTalkLogic {
	return &LikeTalkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeTalkLogic) LikeTalk(req *types.LikeTalkReq) (resp *types.LikeResult, err error) {
	out, err := l.svcCtx.DiscussionService.LikeTalk(l.ctx, &discussionservice.LikeTalkRequest{
		Id: req.TalkId,
	})
	if err != nil {
		return nil, err
	}

	return &types.LikeResult{Liked: out.Liked, Changed: out.Changed}, nil
}
