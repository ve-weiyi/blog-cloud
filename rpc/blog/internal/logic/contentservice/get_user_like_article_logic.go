package contentservicelogic

import (
	"context"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/constants/cachekey"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/contentrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type GetUserLikeArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLikeArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLikeArticleLogic {
	return &GetUserLikeArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询用户点赞的文章
func (l *GetUserLikeArticleLogic) GetUserLikeArticle(in *contentrpc.GetUserLikeArticleRequest) (*contentrpc.GetUserLikeArticleResponse, error) {
	likeKey := cachekey.GetUserLikeArticleKey(in.UserId)
	result, err := l.svcCtx.Redis.SMembers(l.ctx, likeKey).Result()
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0)
	for _, v := range result {
		ids = append(ids, cast.ToInt64(v))
	}

	return &contentrpc.GetUserLikeArticleResponse{
		LikeArticleIds: ids,
	}, nil
}
