package contentservicelogic

import (
	"context"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/contentrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/cachekey"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
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
	uid, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil || uid == "" {
		return nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, "缺少用户身份")
	}

	result, err := l.svcCtx.Redis.SMembers(l.ctx, cachekey.GetUserLikeArticleKey(uid)).Result()
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(result))
	for _, v := range result {
		ids = append(ids, cast.ToInt64(v))
	}

	return &contentrpc.GetUserLikeArticleResponse{
		LikeArticleIds: ids,
	}, nil
}
