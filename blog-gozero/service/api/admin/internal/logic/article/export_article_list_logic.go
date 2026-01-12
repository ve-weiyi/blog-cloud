package article

import (
	"context"
	"path"
	"time"

	"github.com/ve-weiyi/ve-blog-golang/pkg/kit/quickstart/gotplgen"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/articlerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportArticleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 导出文章列表
func NewExportArticleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportArticleListLogic {
	return &ExportArticleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportArticleListLogic) ExportArticleList(req *types.IdsReq) (resp *types.EmptyResp, err error) {
	in := &articlerpc.FindArticleListReq{}

	out, err := l.svcCtx.ArticleRpc.FindArticleList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var list []*types.ArticleBackVO
	for _, v := range out.List {
		m := ConvertArticleTypes(v)
		list = append(list, m)
	}

	// 转换数据
	for _, v := range list {
		err = l.exportArticle(v)
		if err != nil {
			return nil, err
		}
	}

	return &types.EmptyResp{}, nil
}

func (l *ExportArticleListLogic) exportArticle(a *types.ArticleBackVO) (err error) {
	fn := path.Join("./runtime/article", a.ArticleTitle+".md")

	ac := gotplgen.TemplateMeta{
		Mode:           gotplgen.ModeCreateOrReplace,
		CodeOutPath:    fn,
		TemplateString: articleTemplate,
		FunMap:         nil,
		Data: map[string]any{
			"ArticleTitle":    a.ArticleTitle,
			"ArticleCover":    a.ArticleCover,
			"ArticleType":     a.ArticleType,
			"ArticleCategory": a.CategoryName,
			"ArticleTags":     a.TagNameList,
			"ArticleContent":  a.ArticleContent,
			"CreateTime":      time.Unix(a.CreatedAt, 0).String(),
		},
	}

	return ac.Execute()
}

const articleTemplate = `
# {{.ArticleTitle}}  
文章封面: {{.ArticleCover}}   
文章类型: {{.ArticleType}}   
文章分类: {{.ArticleCategory}}   
文章标签: {{.ArticleTags}}   
创建时间: {{.CreateTime}}   

文章内容:
{{.ArticleContent}}
`
