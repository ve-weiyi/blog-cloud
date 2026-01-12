package constant

// 相册相关枚举 (t_album / t_photo)
const (
	// 相册删除状态 (is_delete)
	AlbumIsDeleteNo  = 0 // 未删除
	AlbumIsDeleteYes = 1 // 已删除

	// 照片删除状态 (is_delete)
	PhotoIsDeleteNo  = 0 // 未删除
	PhotoIsDeleteYes = 1 // 已删除

	// 相册状态 (status)
	AlbumStatusPublic  = 1 // 公开
	AlbumStatusPrivate = 2 // 私密
)

// API相关枚举 (t_api)
const (
	// API追溯操作记录 (traceable)
	APITraceableNo  = 0 // 不需要追溯操作记录
	APITraceableYes = 1 // 需要追溯操作记录

	// API禁用状态 (status)
	APIStatusNormal   = 0 // 正常
	APIStatusDisabled = 1 // 禁用
)

// 文章相关枚举 (t_article)
const (
	// 文章类型 (article_type)
	ArticleTypeOriginal  = 1 // 原创
	ArticleTypeReprint   = 2 // 转载
	ArticleTypeTranslate = 3 // 翻译

	// 文章置顶状态 (is_top)
	ArticleIsTopALL = -1
	ArticleIsTopNo  = 0 // 未置顶
	ArticleIsTopYes = 1 // 已置顶

	// 文章删除状态 (is_delete)
	ArticleIsDeleteALL = -1
	ArticleIsDeleteNo  = 0 // 未删除
	ArticleIsDeleteYes = 1 // 已删除

	// 文章状态 (status)
	ArticleStatusPublic         = 1 // 公开
	ArticleStatusPrivate        = 2 // 私密
	ArticleStatusDraft          = 3 // 草稿
	ArticleStatusCommentVisible = 4 // 评论可见
)

// 评论相关枚举 (t_comment)
const (
	// 评论类型 (type)
	CommentTypeArticle = 1 // 文章评论
	CommentTypeFriend  = 2 // 友链评论
	CommentTypeTalk    = 3 // 说说评论

	// 评论状态 (status)
	CommentStatusNormal  = 0 // 正常
	CommentStatusEdited  = 1 // 已编辑
	CommentStatusDeleted = 2 // 已删除

	// 评论审核状态 (is_review)
	CommentIsReviewNo  = 0 // 未审核通过
	CommentIsReviewYes = 1 // 已审核通过
)

// 聊天相关枚举 (t_chat)
const (
	// 聊天消息状态 (status)
	ChatStatusNormal  = 0 // 正常
	ChatStatusEdited  = 1 // 编辑
	ChatStatusRecall  = 2 // 撤回
	ChatStatusDeleted = 3 // 删除
)

// 留言相关枚举 (t_remark)
const (
	// 留言状态 (status)
	RemarkStatusNormal  = 0 // 正常
	RemarkStatusEdited  = 1 // 编辑
	RemarkStatusRecall  = 2 // 撤回
	RemarkStatusDeleted = 3 // 删除

	// 留言审核状态 (is_review)
	RemarkIsReviewNo  = 0 // 未审核通过
	RemarkIsReviewYes = 1 // 已审核通过
)

// 页面相关枚举 (t_page)
const (
	// 页面轮播状态 (is_carousel)
	PageIsCarouselNo  = 0 // 非轮播
	PageIsCarouselYes = 1 // 轮播
)

// 角色相关枚举 (t_role)
const (
	// 角色禁用状态 (status)
	RoleStatusNormal   = 0 // 正常
	RoleStatusDisabled = 1 // 禁用

	// 角色默认状态 (is_default)
	RoleIsDefaultNo  = 0 // 非默认角色
	RoleIsDefaultYes = 1 // 默认角色
)

// 说说相关枚举 (t_talk)
const (
	// 说说置顶状态 (is_top)
	TalkIsTopNo  = 0 // 未置顶
	TalkIsTopYes = 1 // 已置顶

	// 说说状态 (status)
	TalkStatusPublic  = 1 // 公开
	TalkStatusPrivate = 2 // 私密
)

// 用户相关枚举 (t_user)
const (
	// 用户状态 (status)
	UserStatusDeleted  = -1 // 已删除
	UserStatusNormal   = 0  // 正常
	UserStatusDisabled = 1  // 禁用
)

// 登录日志 (t_login_log)
const (
	// 登录类型
	LoginTypeUsername = "username" // 用户名登录
	LoginTypeEmail    = "email"    // 邮箱登录
	LoginTypePhone    = "phone"    // 手机登录
	LoginTypeOauth    = "oauth"    // 第三方登录
)

// 访问统计相关枚举 (t_visit_daily_stats)
const (
	// 访问类型 (visit_type)
	VisitTypeUv = 1 // 独立访客 uv
	VisitTypePv = 2 // 页面浏览量 pv
)

// 菜单相关枚举 (t_menu)
const (
	// 菜单缓存状态 (keep_alive)
	MenuKeepAliveNo  = 0 // 不缓存
	MenuKeepAliveYes = 1 // 缓存

	// 菜单一直显示状态 (always_show)
	MenuAlwaysShowNo  = 0 // 不一直显示
	MenuAlwaysShowYes = 1 // 一直显示

	// 菜单隐藏状态 (is_hidden)
	MenuVisibleNo  = 0 // 隐藏
	MenuVisibleYes = 1 // 显示

	// 菜单禁用状态 (status)
	MenuStatusNormal   = 0 // 正常
	MenuStatusDisabled = 1 // 禁用
)
