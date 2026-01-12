package messagerpclogic

import (
	"context"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/rediskey"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLikeCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLikeCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLikeCommentLogic {
	return &FindUserLikeCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户点赞的评论
func (l *FindUserLikeCommentLogic) FindUserLikeComment(in *messagerpc.UserIdReq) (*messagerpc.FindLikeCommentResp, error) {
	uid := cast.ToString(in.UserId)
	likeKey := rediskey.GetUserLikeCommentKey(uid)

	result, err := l.svcCtx.Redis.SMembers(l.ctx, likeKey).Result()
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0)
	for _, v := range result {
		ids = append(ids, cast.ToInt64(v))
	}

	return &messagerpc.FindLikeCommentResp{
		LikeCommentList: ids,
	}, nil
}
