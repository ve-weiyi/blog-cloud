package articlerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

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
func (l *FindArticlePreviewListLogic) FindArticlePreviewList(in *articlerpc.FindArticlePreviewListReq) (*articlerpc.FindArticlePreviewListResp, error) {
	// 查询文章信息
	records, err := l.svcCtx.TArticleModel.FindALL(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	var list []*articlerpc.ArticlePreview
	for _, v := range records {
		list = append(list, convertArticlePreviewOut(v))
	}

	resp := &articlerpc.FindArticlePreviewListResp{}
	resp.List = list
	return resp, nil
}
