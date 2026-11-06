package comment

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/discussionservice"
)

type UnlikeCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 取消点赞评论
func NewUnlikeCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnlikeCommentLogic {
	return &UnlikeCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnlikeCommentLogic) UnlikeComment(req *types.LikeCommentReq) (resp *types.LikeResult, err error) {
	out, err := l.svcCtx.DiscussionService.UnlikeComment(l.ctx, &discussionservice.LikeCommentRequest{
		Id: req.CommentId,
	})
	if err != nil {
		return nil, err
	}

	return &types.LikeResult{Liked: out.Liked, Changed: out.Changed}, nil
}
