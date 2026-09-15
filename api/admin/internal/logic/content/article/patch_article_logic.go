package article

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/contentservice"
)

type PatchArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新文章删除状态
func NewPatchArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PatchArticleLogic {
	return &PatchArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PatchArticleLogic) PatchArticle(req *types.PatchArticleReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.ContentService.PatchArticle(l.ctx, &contentservice.PatchArticleRequest{
		Id:       req.Id,
		IsDelete: req.IsDelete,
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
