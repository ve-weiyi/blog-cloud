package article

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/articlerpc"

	"github.com/zeromicro/go-zero/core/logx"
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
	in := &articlerpc.UpdateArticleDeleteReq{
		ArticleId: req.Id,
		IsDelete:  req.IsDelete,
	}

	_, err = l.svcCtx.ArticleRpc.UpdateArticleDelete(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
