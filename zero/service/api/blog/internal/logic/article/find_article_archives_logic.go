package article

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/articlerpc"
)

type FindArticleArchivesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文章归档(时间轴)
func NewFindArticleArchivesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindArticleArchivesLogic {
	return &FindArticleArchivesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindArticleArchivesLogic) FindArticleArchives(req *types.ArticleArchivesQueryReq) (resp *types.PageResp, err error) {
	in := &articlerpc.FindArticleListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Sorts:    req.Sorts,
	}
	out, err := l.svcCtx.ArticleRpc.FindArticlePublicList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	list := make([]*types.ArticleHome, 0)
	// 转换数据
	for _, v := range out.List {
		m := ConvertArticleHomeTypes(v)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	resp.Total = out.Total
	resp.List = list
	return
}
