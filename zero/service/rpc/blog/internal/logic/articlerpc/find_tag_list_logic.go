package articlerpclogic

import (
	"context"
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindTagListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindTagListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindTagListLogic {
	return &FindTagListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询标签数量
func (l *FindTagListLogic) FindTagList(in *articlerpc.FindTagListReq) (*articlerpc.FindTagListResp, error) {
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

	if in.TagName != "" {
		conditions += "tag_name like ?"
		params = append(params, "%"+in.TagName+"%")
	}

	records, err := l.svcCtx.TTagModel.FindList(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	count, err := l.svcCtx.TTagModel.FindCount(l.ctx, conditions, params...)
	if err != nil {
		return nil, err
	}

	acm, err := helper.findArticleCountGroupTag(records)
	if err != nil {
		return nil, err
	}

	var list []*articlerpc.TagDetails
	for _, v := range records {
		list = append(list, convertTagOut(v, acm))
	}

	return &articlerpc.FindTagListResp{
		List:  list,
		Total: count,
	}, nil
}
