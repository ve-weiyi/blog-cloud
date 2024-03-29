package response

// 博客前台信息
type BlogHomeInfo struct {
	ArticleCount  int64            `json:"article_count"`  // 文章数量
	CategoryCount int64            `json:"category_count"` // 分类数量
	TagCount      int64            `json:"tag_count"`      // 标签数量
	ViewsCount    string           `json:"views_count"`    // 访问量
	WebsiteConfig WebsiteConfigDTO `json:"website_config"` // 网站配置
	PageList      []*PageDTO       `json:"page_list"`      // 页面列表
}

type WebsiteConfigDTO struct {
	AdminUrl          string      `json:"admin_url"`           // 后台地址
	AlipayQrCode      string      `json:"alipay_qr_code"`      // 支付宝二维码
	Gitee             string      `json:"gitee"`               // Gitee
	Github            string      `json:"github"`              // Github
	IsChatRoom        int         `json:"is_chat_room"`        // 是否开启聊天室
	IsCommentReview   int         `json:"is_comment_review"`   // 是否开启评论审核
	IsEmailNotice     int         `json:"is_email_notice"`     // 是否开启邮件通知
	IsMessageReview   int         `json:"is_message_review"`   // 是否开启留言审核
	IsMusicPlayer     int         `json:"is_music_player"`     // 是否开启音乐播放器
	IsReward          int         `json:"is_reward"`           // 是否开启打赏
	Qq                string      `json:"qq"`                  // QQ
	SocialLoginList   []string    `json:"social_login_list"`   // 社交登录列表
	SocialUrlList     []string    `json:"social_url_list"`     // 社交地址列表
	TouristAvatar     string      `json:"tourist_avatar"`      // 游客头像
	UserAvatar        string      `json:"user_avatar"`         // 用户头像
	WebsiteAuthor     string      `json:"website_author"`      // 网站作者
	WebsiteAvatar     interface{} `json:"website_avatar"`      // 网站头像
	WebsiteCreateTime string      `json:"website_create_time"` // 网站创建时间
	WebsiteIntro      string      `json:"website_intro"`       // 网站介绍
	WebsiteName       string      `json:"website_name"`        // 网站名称
	WebsiteNotice     string      `json:"website_notice"`      // 网站公告
	WebsiteRecordNo   string      `json:"website_record_no"`   // 网站备案号
	WebsocketUrl      string      `json:"websocket_url"`       // websocket地址
	WeixinQrCode      string      `json:"weixin_qr_code"`      // 微信二维码
}

type PageDTO struct {
	ID        int    `json:"id"`         // 页面ID
	PageName  string `json:"page_name"`  // 页面名称
	PageLabel string `json:"page_label"` // 页面标签
	PageCover string `json:"page_cover"` // 页面封面
}

// 博客后台信息
type AdminHomeInfo struct {
	ViewsCount            int64                   `json:"views_count"`             // 访问量
	MessageCount          int64                   `json:"message_count"`           // 留言量
	UserCount             int64                   `json:"user_count"`              // 用户量
	ArticleCount          int64                   `json:"article_count"`           // 文章量
	CategoryDTOList       []*CategoryDTO          `json:"category_dto_list"`       // 分类统计
	TagDTOList            []*TagDTO               `json:"tag_dto_list"`            // 标签列表
	ArticleStatisticsList []*ArticleStatisticsDTO `json:"article_statistics_list"` // 文章统计列表
	UniqueViewDTOList     []*UniqueViewDTO        `json:"unique_view_dto_list"`    // 一周用户量集合
	ArticleRankDTOList    []*ArticleRankDTO       `json:"article_rank_dto_list"`   // 文章浏览量排行
}

// 文章统计字段
type ArticleStatisticsDTO struct {
	Day   string `json:"day"`   // 日期
	Count int64  `json:"count"` // 数量
}

// 一周用户量
type UniqueViewDTO struct {
	Day   string `json:"day"`   // 日期
	Count int64  `json:"count"` // 数量
}

// 文章浏览量排行
type ArticleRankDTO struct {
	ID           int    `json:"id"`            // 文章ID
	ArticleTitle string `json:"article_title"` // 文章标题
	Count        int64  `json:"count"`         // 数量
}
