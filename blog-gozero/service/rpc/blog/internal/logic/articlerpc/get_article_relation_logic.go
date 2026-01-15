package articlerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/constant"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArticleRelationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetArticleRelationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleRelationLogic {
	return &GetArticleRelationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询关联文章
func (l *GetArticleRelationLogic) GetArticleRelation(in *articlerpc.GetArticleRelationReq) (*articlerpc.GetArticleRelationResp, error) {

	record, err := l.svcCtx.TArticleModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 查询上一篇文章
	last, _, err := l.svcCtx.TArticleModel.FindListAndTotal(l.ctx, 1, 1, "id desc", "id < ? and status = ?", record.Id, constant.ArticleStatusPublic)
	if err != nil {
		return nil, err
	}
	// 查询下一篇文章
	next, _, err := l.svcCtx.TArticleModel.FindListAndTotal(l.ctx, 1, 1, "id asc", "id > ? and status = ?", record.Id, constant.ArticleStatusPublic)
	if err != nil {
		return nil, err
	}
	// 查询推荐文章
	recommend, _, err := l.svcCtx.TArticleModel.FindListAndTotal(l.ctx, 1, 5, "id asc", "id != ? and category_id = ? and status = ?", record.Id, record.CategoryId, constant.ArticleStatusPublic)
	if err != nil {
		return nil, err
	}
	// 查询最新文章
	newest, _, err := l.svcCtx.TArticleModel.FindListAndTotal(l.ctx, 1, 5, "id desc", "status = ?", constant.ArticleStatusPublic)
	if err != nil {
		return nil, err
	}

	helper := NewArticleHelperLogic(l.ctx, l.svcCtx)
	var la, na *articlerpc.ArticlePreview
	var ras, nas []*articlerpc.ArticlePreview
	if len(last) > 0 {
		la = helper.convertArticlePreviewOut(last[0])
	}

	if len(next) > 0 {
		na = helper.convertArticlePreviewOut(next[0])
	}

	if len(recommend) > 0 {
		for _, v := range recommend {
			ras = append(ras, helper.convertArticlePreviewOut(v))
		}
	}

	if len(newest) > 0 {
		for _, v := range newest {
			nas = append(nas, helper.convertArticlePreviewOut(v))
		}
	}

	resp := &articlerpc.GetArticleRelationResp{
		Last:      la,
		Next:      na,
		Recommend: ras,
		Newest:    nas,
	}

	return resp, nil
}
