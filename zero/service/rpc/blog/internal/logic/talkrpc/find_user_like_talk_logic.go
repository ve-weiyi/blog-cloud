package talkrpclogic

import (
	"context"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/common/rediskey"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/talkrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLikeTalkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLikeTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLikeTalkLogic {
	return &FindUserLikeTalkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户点赞的说说
func (l *FindUserLikeTalkLogic) FindUserLikeTalk(in *talkrpc.UserIdReq) (*talkrpc.FindLikeTalkResp, error) {
	uid := cast.ToString(in.UserId)
	likeKey := rediskey.GetUserLikeTalkKey(uid)

	result, err := l.svcCtx.Redis.SMembers(l.ctx, likeKey).Result()
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0)
	for k, _ := range result {
		ids = append(ids, cast.ToInt64(k))
	}

	return &talkrpc.FindLikeTalkResp{
		LikeTalkList: ids,
	}, nil
}
