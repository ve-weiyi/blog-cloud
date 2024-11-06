package articlerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticleRecommendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetArticleRecommendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleRecommendLogic {
	return &GetArticleRecommendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询文章推荐
func (l *GetArticleRecommendLogic) GetArticleRecommend(in *articlerpc.IdReq) (*articlerpc.ArticleRecommendResp, error) {

	record, err := l.svcCtx.TArticleModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 查询上一篇文章
	last, err := l.svcCtx.TArticleModel.FindList(l.ctx, 1, 1, "id desc", "id < ?", record.Id)
	if err != nil {
		return nil, err
	}
	// 查询下一篇文章
	next, err := l.svcCtx.TArticleModel.FindList(l.ctx, 1, 1, "id asc", "id > ?", record.Id)
	if err != nil {
		return nil, err
	}
	// 查询推荐文章
	recommend, err := l.svcCtx.TArticleModel.FindList(l.ctx, 1, 5, "id asc", "id != ? and category_id = ?", record.Id, record.CategoryId)
	if err != nil {
		return nil, err
	}
	// 查询最新文章
	newest, err := l.svcCtx.TArticleModel.FindList(l.ctx, 1, 5, "id desc", "")
	if err != nil {
		return nil, err
	}

	var la, na *articlerpc.ArticlePreview
	var ras, nas []*articlerpc.ArticlePreview
	if len(last) > 0 {
		la = convertArticlePreviewOut(last[0])
	}

	if len(next) > 0 {
		na = convertArticlePreviewOut(next[0])
	}

	if len(recommend) > 0 {
		for _, v := range recommend {
			ras = append(ras, convertArticlePreviewOut(v))
		}
	}

	if len(newest) > 0 {
		for _, v := range newest {
			nas = append(nas, convertArticlePreviewOut(v))
		}
	}

	resp := &articlerpc.ArticleRecommendResp{
		Last:      la,
		Next:      na,
		Recommend: ras,
		Newest:    nas,
	}

	return resp, nil
}
