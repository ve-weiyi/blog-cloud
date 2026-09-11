package cachekey

import "fmt"

// 命名规则：服务:模块:操作:参数

// ──────────── 有参函数 ────────────

// 验证码 redis类型：string
func GetCaptchaKey(module string, username string) string {
	return fmt.Sprintf("blog:captcha:%s:%s", module, username)
}

// 用户点赞的文章集合 redis类型：set
func GetUserLikeArticleKey(uid string) string {
	return fmt.Sprintf("blog:user:like:article:%s", uid)
}

// 用户点赞的评论集合 redis类型：set
func GetUserLikeCommentKey(uid string) string {
	return fmt.Sprintf("blog:user:like:comment:%s", uid)
}

// 用户点赞的说说集合 redis类型：set
func GetUserLikeTalkKey(uid string) string {
	return fmt.Sprintf("blog:user:like:talk:%s", uid)
}

// 网站日访客集合
func GetDailyUserVisitKey(day string) string {
	return fmt.Sprintf("blog:visit:visitor:%v", day)
}

// ──────────── 无参常量 ────────────

const (
	// 文章点赞数排行
	ArticleLikeCountKey = "blog:article:like_count"
	// 评论点赞数排行
	CommentLikeCountKey = "blog:comment:like_count"
	// 说说点赞数排行
	TalkLikeCountKey = "blog:talk:like_count"
	// 文章访问量排行
	ArticleViewCountKey = "blog:article:view_count"
	// 网站日 UV
	DailyUserViewCountKey = "blog:visit:daily_uv"
	// 网站日 PV
	DailyPageViewCountKey = "blog:visit:daily_pv"
	// 总 UV
	TotalUserViewCountKey = "blog:visit:total_uv"
	// 总 PV
	TotalPageViewCountKey = "blog:visit:total_pv"
	// 在线用户集合
	OnlineUserKey = "blog:online:user"
	// 在线管理员集合
	OnlineAdminKey = "blog:online:admin"
)

const (
	// Token 存储前缀
	TokenStorePrefixApp   = "blog:app:token:"
	TokenStorePrefixAdmin = "blog:admin:token:"
	// 验证码存储前缀
	CaptchaStorePrefixApp   = "blog:app:captcha:"
	CaptchaStorePrefixAdmin = "blog:admin:captcha:"
	// 限流器前缀
	RateLimitStrictPrefix = "blog:ratelimit:strict:"
)
