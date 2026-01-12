package articlerpclogic

import (
	"context"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/rediskey"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddArticleVisitsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddArticleVisitsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddArticleVisitsLogic {
	return &AddArticleVisitsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 访问文章
func (l *AddArticleVisitsLogic) AddArticleVisits(in *articlerpc.AddArticleVisitsReq) (*articlerpc.AddArticleVisitsResp, error) {
	record, err := l.svcCtx.TArticleModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	record.ViewCount++
	_, err = l.svcCtx.TArticleModel.Updates(l.ctx, map[string]interface{}{
		"view_count": record.ViewCount,
	},
		"id = ?", record.Id,
	)
	if err != nil {
		return nil, err
	}

	id := cast.ToString(record.Id)
	key := rediskey.GetArticleViewCountKey()
	// 浏览量+1
	_, err = l.svcCtx.Redis.ZIncrBy(l.ctx, key, 1, id).Result()
	if err != nil {
		return nil, err
	}

	return &articlerpc.AddArticleVisitsResp{}, nil
}
