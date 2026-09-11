package article

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/articleservice"
)

type UpdateArticleTopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新文章置顶状态
func NewUpdateArticleTopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateArticleTopLogic {
	return &UpdateArticleTopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateArticleTopLogic) UpdateArticleTop(req *types.UpdateArticleTopReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.ArticleService.PatchArticle(l.ctx, &articleservice.PatchArticleRequest{
		Id:    req.Id,
		IsTop: &req.IsTop,
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
