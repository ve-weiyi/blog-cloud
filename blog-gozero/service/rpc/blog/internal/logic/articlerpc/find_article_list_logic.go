package articlerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindArticleListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindArticleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindArticleListLogic {
	return &FindArticleListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询文章数量
func (l *FindArticleListLogic) FindArticleList(in *articlerpc.FindArticleListReq) (*articlerpc.FindArticleListResp, error) {
	helper := NewArticleHelperLogic(l.ctx, l.svcCtx)

	page, size, sorts, conditions, params := helper.convertArticleQuery(in)

	// 查询文章信息
	records, total, err := l.svcCtx.TArticleModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	list, err := helper.convertArticleDetailsResp(records)
	if err != nil {
		return nil, err
	}

	return &articlerpc.FindArticleListResp{
		List: list,
		Pagination: &articlerpc.PageResp{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
	}, nil
}
