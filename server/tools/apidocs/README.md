## 前端 api.ts 生成器


### 1. 介绍

根据swagger.json 或者 handler 注释，解析出 api.ts 文件，提供给前端使用。
可以同步后端接口修改影响，减少前后端接口对接成本。


### 2. 使用方式

1. 在controller接口上添加注释，如
```golang
// @Tags 	 	Article
// @Summary		分页获取文章列表
// @Accept 		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param 	 	page 	body		request.PageQuery 			true 	"分页参数"
// @Success		200		{object}	response.Response{data=response.PageResult{list=[]response.ArticleBack}}	"返回信息"
// @Router		/admin/article/list [post]
func (s *ArticleController) FindArticleList(c *gin.Context) {
    // ...
}
```

2. 运行test方法,生成的ts文件自动保存在[api](api)

```golang
func TestApiDocs(t *testing.T) {
	cfg := Config{
        // ...
    }

	aad := NewAstApiDoc(cfg)
	aad.Parse()
	// 生成ts api定义文件
	aad.GenerateTsTypeFile()
	// 生成ts type定义文件
	aad.GenerateTsApiFiles()
}

```

生成结果如下

article.ts
```ts

/** 分页获取文章列表 */
export function findArticleListApi(page: PageQuery): Promise<IApiResponseData<PageResult<ArticleBack[]>>> {
    return http.request<IApiResponseData<PageResult<ArticleBack[]>>>({
        url: `/api/v1/admin/article/list`,
        method: "post",
        data: page,
    })
}
```

type.ts
```ts
export interface ArticleBack extends ArticleDTO {
    category_name?: string // 文章分类名
    tag_name_list?: string[] // 文章标签列表
}
```
3. 移动api目录到前端项目中

