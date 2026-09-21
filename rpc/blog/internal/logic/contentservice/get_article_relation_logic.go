package contentservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/constants/enums"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/contentrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
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
func (l *GetArticleRelationLogic) GetArticleRelation(in *contentrpc.GetArticleRelationRequest) (*contentrpc.GetArticleRelationResponse, error) {
	helper := NewArticleHelper(l.ctx, l.svcCtx)

	record, err := l.svcCtx.TArticleModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	last, _, err := l.svcCtx.TArticleModel.FindListAndTotal(l.ctx, 1, 1, "id desc", "id < ? and status = ?", record.Id, enums.ArticleStatusPublic)
	if err != nil {
		l.Errorf("GetArticleRelation FindListAndTotal last error: %v", err)
	}
	next, _, err := l.svcCtx.TArticleModel.FindListAndTotal(l.ctx, 1, 1, "id asc", "id > ? and status = ?", record.Id, enums.ArticleStatusPublic)
	if err != nil {
		l.Errorf("GetArticleRelation FindListAndTotal next error: %v", err)
	}
	recommend, _, err := l.svcCtx.TArticleModel.FindListAndTotal(l.ctx, 1, 5, "id asc", "id != ? and category_id = ? and status = ?", record.Id, record.CategoryId, enums.ArticleStatusPublic)
	if err != nil {
		l.Errorf("GetArticleRelation FindListAndTotal recommend error: %v", err)
	}
	newest, _, err := l.svcCtx.TArticleModel.FindListAndTotal(l.ctx, 1, 5, "id desc", "status = ?", enums.ArticleStatusPublic)
	if err != nil {
		l.Errorf("GetArticleRelation FindListAndTotal newest error: %v", err)
	}

	resp := &contentrpc.GetArticleRelationResponse{}
	if len(last) > 0 {
		resp.Last = helper.convertArticlePreviewOut(last[0])
	}
	if len(next) > 0 {
		resp.Next = helper.convertArticlePreviewOut(next[0])
	}
	for _, v := range recommend {
		resp.Recommends = append(resp.Recommends, helper.convertArticlePreviewOut(v))
	}
	for _, v := range newest {
		resp.Latest = append(resp.Latest, helper.convertArticlePreviewOut(v))
	}

	return resp, nil
}
