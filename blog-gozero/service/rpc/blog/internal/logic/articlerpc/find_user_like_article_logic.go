package articlerpclogic

import (
	"context"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/rediskey"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLikeArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLikeArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLikeArticleLogic {
	return &FindUserLikeArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户点赞的文章
func (l *FindUserLikeArticleLogic) FindUserLikeArticle(in *articlerpc.UserIdReq) (*articlerpc.FindLikeArticleResp, error) {
	uid := cast.ToString(in.UserId)
	likeKey := rediskey.GetUserLikeArticleKey(uid)

	result, err := l.svcCtx.Redis.SMembers(l.ctx, likeKey).Result()
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0)
	for _, v := range result {
		ids = append(ids, cast.ToInt64(v))
	}

	return &articlerpc.FindLikeArticleResp{
		LikeArticleList: ids,
	}, nil
}
