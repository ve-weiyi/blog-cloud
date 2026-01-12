package articlerpclogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
)

type UpdateArticleTopLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateArticleTopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateArticleTopLogic {
	return &UpdateArticleTopLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新文章置顶
func (l *UpdateArticleTopLogic) UpdateArticleTop(in *articlerpc.UpdateArticleTopReq) (*articlerpc.UpdateArticleTopResp, error) {

	record, err := l.svcCtx.TArticleModel.FindById(l.ctx, in.ArticleId)
	if err != nil {
		return nil, err
	}

	record.IsTop = in.IsTop
	_, err = l.svcCtx.TArticleModel.Save(l.ctx, record)
	if err != nil {
		return nil, err
	}

	helper := NewArticleHelperLogic(l.ctx, l.svcCtx)
	return &articlerpc.UpdateArticleTopResp{
		Article: helper.convertArticlePreviewOut(record),
	}, nil
}
