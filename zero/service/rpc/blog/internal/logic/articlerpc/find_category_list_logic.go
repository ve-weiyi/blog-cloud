package articlerpclogic

import (
	"context"
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindCategoryListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindCategoryListLogic {
	return &FindCategoryListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询文章分类数量
func (l *FindCategoryListLogic) FindCategoryList(in *articlerpc.FindCategoryListReq) (*articlerpc.FindCategoryListResp, error) {
	helper := NewArticleHelperLogic(l.ctx, l.svcCtx)
	var (
		page       int
		size       int
		sorts      string
		conditions string
		params     []interface{}
	)

	page = int(in.Page)
	size = int(in.PageSize)
	sorts = strings.Join(in.Sorts, ",")
	if sorts == "" {
		sorts = "id desc"
	}
	if in.CategoryName != "" {
		conditions += "category_name like ?"
		params = append(params, "%"+in.CategoryName+"%")
	}

	records, err := l.svcCtx.TCategoryModel.FindList(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	count, err := l.svcCtx.TCategoryModel.FindCount(l.ctx, conditions, params...)
	if err != nil {
		return nil, err
	}

	acm, err := helper.findArticleCountGroupCategory(records)
	if err != nil {
		return nil, err
	}

	var list []*articlerpc.CategoryDetails
	for _, v := range records {
		list = append(list, convertCategoryOut(v, acm))
	}

	return &articlerpc.FindCategoryListResp{
		List:  list,
		Total: count,
	}, nil
}
