package contentservicelogic

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/contentrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/cachekey"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
)

type BatchDeleteArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteArticlesLogic {
	return &BatchDeleteArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除文章
func (l *BatchDeleteArticlesLogic) BatchDeleteArticles(in *contentrpc.BatchDeleteArticlesRequest) (*contentrpc.BatchDeleteArticlesResponse, error) {
	var rows int64
	err := l.svcCtx.GormDB.Transaction(func(tx *gorm.DB) error {
		if _, err := l.svcCtx.TArticleTagModel.WithTx(tx).DeleteBatch(l.ctx, "article_id in (?)", in.Ids); err != nil {
			return err
		}

		if _, err := l.svcCtx.TLikeModel.WithTx(tx).DeleteBatch(l.ctx, "target_type = ? AND target_id in (?)", enums.LikeTargetArticle, in.Ids); err != nil {
			return err
		}

		var err error
		rows, err = l.svcCtx.TArticleModel.WithTx(tx).DeleteBatch(l.ctx, "id in (?)", in.Ids)
		return err
	})
	if err != nil {
		return nil, err
	}

	// 计数缓存是派生数据，清理失败只记录日志
	members := make([]interface{}, 0, len(in.Ids))
	for _, id := range in.Ids {
		members = append(members, fmt.Sprintf("%d", id))
	}
	if err := l.svcCtx.Redis.ZRem(l.ctx, cachekey.ArticleLikeCountKey, members...).Err(); err != nil {
		l.Errorf("清理文章点赞计数缓存失败: %v", err)
	}
	if err := l.svcCtx.Redis.ZRem(l.ctx, cachekey.ArticleViewCountKey, members...).Err(); err != nil {
		l.Errorf("清理文章浏览量缓存失败: %v", err)
	}

	return &contentrpc.BatchDeleteArticlesResponse{
		SuccessCount: rows,
	}, nil
}
