package visitx

type Enforcer interface {
	IsVisitPage(path string, method string) (pageName string, ok bool)
}

type defaultEnforcer struct {
	pageMap map[string]string
}

func NewVisitEnforcer() Enforcer {
	return &defaultEnforcer{
		pageMap: map[string]string{
			"GET /api/v1/blog":                          "首页",
			"GET /api/v1/blog/get_about_me":             "关于我",
			"GET /api/v1/album/get_album":               "相册详情",
			"GET /api/v1/article/get_article":           "文章详情",
			"GET /api/v1/talk/get_talk":                 "说说详情",
			"POST /api/v1/message/query_message_list":   "留言",
			"POST /api/v1/category/query_category_list": "分类",
			"POST /api/v1/tag/query_tag_list":           "标签",
		},
	}
}

func (e *defaultEnforcer) IsVisitPage(path string, method string) (string, bool) {
	page, ok := e.pageMap[method+" "+path]
	return page, ok
}
