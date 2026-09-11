package articleservicelogic

import (
	"context"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/constants/cachekey"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type IncrementArticleViewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIncrementArticleViewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IncrementArticleViewLogic {
	return &IncrementArticleViewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 增加文章浏览量
func (l *IncrementArticleViewLogic) IncrementArticleView(in *articlerpc.IncrementArticleViewRequest) (*articlerpc.IncrementArticleViewResponse, error) {
	record, err := l.svcCtx.TArticleModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	record.ViewCount++
	_, err = l.svcCtx.TArticleModel.UpdateFields(l.ctx, map[string]interface{}{
		"view_count": record.ViewCount,
	}, "id = ?", record.Id)
	if err != nil {
		return nil, err
	}

	id := cast.ToString(record.Id)
	key := cachekey.ArticleViewCountKey
	_, err = l.svcCtx.Redis.ZIncrBy(l.ctx, key, 1, id).Result()
	if err != nil {
		return nil, err
	}

	return &articlerpc.IncrementArticleViewResponse{Success: true}, nil
}
