package model

const (
	//1公开 2私密 3评论可见
	ArticleStatusPublic  = 1
	ArticleStatusPrivate = 2
	ArticleStatusComment = 3
)

const (
	//0未删除 1已删除
	ArticleIsDeleteUnused = -1
	ArticleIsDeleteNo     = 0
	ArticleIsDeleteYes    = 1

	//0未置顶 1置顶
	ArticleIsTopUnused = -1
	ArticleIsTopNo     = 0
	ArticleIsTopYes    = 1
)
