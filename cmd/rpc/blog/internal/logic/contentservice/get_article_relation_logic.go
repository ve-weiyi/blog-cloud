package contentservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/contentrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/model"
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

	last, _, err := l.svcCtx.TArticleModel.FindListAndTotal(l.ctx, 1, 1, "id desc", "id < ? and status = ? and is_delete = 0", record.Id, enums.ArticleStatusPublic)
	if err != nil {
		return nil, err
	}
	next, _, err := l.svcCtx.TArticleModel.FindListAndTotal(l.ctx, 1, 1, "id asc", "id > ? and status = ? and is_delete = 0", record.Id, enums.ArticleStatusPublic)
	if err != nil {
		return nil, err
	}
	recommend, _, err := l.svcCtx.TArticleModel.FindListAndTotal(l.ctx, 1, 5, "id asc", "id != ? and category_id = ? and status = ? and is_delete = 0", record.Id, record.CategoryId, enums.ArticleStatusPublic)
	if err != nil {
		return nil, err
	}
	newest, _, err := l.svcCtx.TArticleModel.FindListAndTotal(l.ctx, 1, 5, "id desc", "status = ? and is_delete = 0", enums.ArticleStatusPublic)
	if err != nil {
		return nil, err
	}

	// 四组预览一次性转换：浏览量走同一条 pipeline，点赞状态只查一次
	all := make([]*model.TArticle, 0, len(last)+len(next)+len(recommend)+len(newest))
	all = append(all, last...)
	all = append(all, next...)
	all = append(all, recommend...)
	all = append(all, newest...)
	previews, err := helper.convertArticlePreviewOuts(all)
	if err != nil {
		return nil, err
	}
	pm := make(map[int64]*contentrpc.ArticlePreview, len(previews))
	for _, v := range previews {
		pm[v.Id] = v
	}

	resp := &contentrpc.GetArticleRelationResponse{}
	if len(last) > 0 {
		resp.Last = pm[last[0].Id]
	}
	if len(next) > 0 {
		resp.Next = pm[next[0].Id]
	}
	for _, v := range recommend {
		resp.Recommends = append(resp.Recommends, pm[v.Id])
	}
	for _, v := range newest {
		resp.Latest = append(resp.Latest, pm[v.Id])
	}

	return resp, nil
}
