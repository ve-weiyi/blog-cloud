package article

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/articleservice"
)

type UpdateArticleDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新文章删除状态
func NewUpdateArticleDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateArticleDeleteLogic {
	return &UpdateArticleDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateArticleDeleteLogic) UpdateArticleDelete(req *types.UpdateArticleDeleteReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.ArticleService.PatchArticle(l.ctx, &articleservice.PatchArticleRequest{
		Id:       req.Id,
		IsDelete: &req.IsDelete,
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
