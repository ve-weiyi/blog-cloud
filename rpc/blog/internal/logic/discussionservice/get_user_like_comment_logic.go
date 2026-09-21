package discussionservicelogic

import (
	"context"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/constants/cachekey"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type GetUserLikeCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLikeCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLikeCommentLogic {
	return &GetUserLikeCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询用户点赞的评论
func (l *GetUserLikeCommentLogic) GetUserLikeComment(in *discussionrpc.GetUserLikeCommentRequest) (*discussionrpc.GetUserLikeCommentResponse, error) {
	likeKey := cachekey.GetUserLikeCommentKey(in.UserId)
	result, err := l.svcCtx.Redis.SMembers(l.ctx, likeKey).Result()
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0)
	for _, v := range result {
		ids = append(ids, cast.ToInt64(v))
	}

	return &discussionrpc.GetUserLikeCommentResponse{LikeCommentIds: ids}, nil
}
