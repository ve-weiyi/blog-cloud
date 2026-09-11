package articleservicelogic

import (
	"context"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/constants/cachekey"
	"github.com/ve-weiyi/blog-cloud/infra/metax"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type LikeArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLikeArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeArticleLogic {
	return &LikeArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 点赞文章
func (l *LikeArticleLogic) LikeArticle(in *articlerpc.LikeArticleRequest) (*articlerpc.LikeArticleResponse, error) {
	uid, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	id := cast.ToString(in.Id)

	entity, err := l.svcCtx.TArticleModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	likeKey := cachekey.GetUserLikeArticleKey(uid)
	countKey := cachekey.ArticleLikeCountKey

	ok, err := l.svcCtx.Redis.SIsMember(l.ctx, likeKey, id).Result()
	if err != nil {
		l.Errorf("LikeArticle SIsMember error: %v", err)
	}
	if ok {
		entity.LikeCount--
		if err := l.svcCtx.Redis.SRem(l.ctx, likeKey, id).Err(); err != nil {
			l.Errorf("LikeArticle SRem error: %v", err)
		}
		if err := l.svcCtx.Redis.ZIncrBy(l.ctx, countKey, -1, id).Err(); err != nil {
			l.Errorf("LikeArticle ZIncrBy error: %v", err)
		}
	} else {
		entity.LikeCount++
		if err := l.svcCtx.Redis.SAdd(l.ctx, likeKey, id).Err(); err != nil {
			l.Errorf("LikeArticle SAdd error: %v", err)
		}
		if err := l.svcCtx.Redis.ZIncrBy(l.ctx, countKey, 1, id).Err(); err != nil {
			l.Errorf("LikeArticle ZIncrBy error: %v", err)
		}
	}

	_, err = l.svcCtx.TArticleModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &articlerpc.LikeArticleResponse{Success: true}, nil
}
