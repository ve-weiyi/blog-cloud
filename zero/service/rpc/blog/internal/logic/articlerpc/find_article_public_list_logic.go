package articlerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindArticlePublicListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindArticlePublicListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindArticlePublicListLogic {
	return &FindArticlePublicListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询文章列表
func (l *FindArticlePublicListLogic) FindArticlePublicList(in *articlerpc.FindArticleListReq) (*articlerpc.FindArticleListResp, error) {
	helper := NewArticleHelperLogic(l.ctx, l.svcCtx)

	in.Status = model.ArticleStatusPublic
	in.IsDelete = model.ArticleIsDeleteNo
	page, size, sorts, conditions, params := helper.convertArticleQuery(in)

	// 查询文章信息
	records, err := l.svcCtx.TArticleModel.FindList(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	count, err := l.svcCtx.TArticleModel.FindCount(l.ctx, conditions, params...)
	if err != nil {
		return nil, err
	}

	acm, err := helper.findCategoryGroupArticle(records)
	if err != nil {
		return nil, err

	}

	atm, err := helper.findTagGroupArticle(records)
	if err != nil {
		return nil, err
	}

	var list []*articlerpc.ArticleDetails
	for _, v := range records {
		list = append(list, convertArticleOut(v, acm, atm))
	}

	return &articlerpc.FindArticleListResp{
		List:  list,
		Total: count,
	}, nil
}
