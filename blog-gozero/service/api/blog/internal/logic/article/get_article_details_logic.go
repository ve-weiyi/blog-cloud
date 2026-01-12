package article

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/common/apiutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"
)

type GetArticleDetailsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取文章详情
func NewGetArticleDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleDetailsLogic {
	return &GetArticleDetailsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticleDetailsLogic) GetArticleDetails(req *types.IdReq) (resp *types.ArticleDetails, err error) {
	// 添加文章访问量
	_, err = l.svcCtx.ArticleRpc.AddArticleVisits(l.ctx, &articlerpc.AddArticleVisitsReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	out, err := l.svcCtx.ArticleRpc.GetArticle(l.ctx, &articlerpc.GetArticleReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	// 查询关联文章
	relation, err := l.svcCtx.ArticleRpc.GetArticleRelation(l.ctx, &articlerpc.GetArticleRelationReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	infos, err := apiutils.GetUserInfos(l.ctx, l.svcCtx, []string{out.Article.UserId})
	if err != nil {
		return nil, err
	}

	resp = &types.ArticleDetails{
		ArticleHome:          types.ArticleHome{},
		Author:               nil,
		LastArticle:          nil,
		NextArticle:          nil,
		RecommendArticleList: nil,
		NewestArticleList:    nil,
	}

	resp.Author = infos[out.Article.UserId]

	resp.ArticleHome = *convertArticleHomeTypes(out.Article)

	resp.LastArticle = convertArticlePreviewTypes(relation.Last)

	resp.NextArticle = convertArticlePreviewTypes(relation.Next)

	for _, v := range relation.Recommend {
		resp.RecommendArticleList = append(resp.RecommendArticleList, convertArticlePreviewTypes(v))
	}

	for _, v := range relation.Newest {
		resp.NewestArticleList = append(resp.NewestArticleList, convertArticlePreviewTypes(v))
	}

	_, err = l.svcCtx.SyslogRpc.AddVisitLog(l.ctx, &syslogrpc.AddVisitLogReq{
		PageName: "文章详情",
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
