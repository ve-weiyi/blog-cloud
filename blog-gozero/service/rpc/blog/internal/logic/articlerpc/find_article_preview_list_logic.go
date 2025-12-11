package articlerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindArticlePreviewListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindArticlePreviewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindArticlePreviewListLogic {
	return &FindArticlePreviewListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询文章列表
func (l *FindArticlePreviewListLogic) FindArticlePreviewList(in *articlerpc.FindArticleListReq) (*articlerpc.FindArticlePreviewListResp, error) {
	helper := NewArticleHelperLogic(l.ctx, l.svcCtx)

	page, size, sorts, conditions, params := helper.convertArticleQuery(in)

	// 查询文章信息
	records, total, err := l.svcCtx.TArticleModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*articlerpc.ArticlePreview
	for _, v := range records {
		list = append(list, helper.convertArticlePreviewOut(v))
	}

	return &articlerpc.FindArticlePreviewListResp{
		List: list,
		Pagination: &articlerpc.PageResp{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
	}, nil
}
