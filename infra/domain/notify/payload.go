package notify

// 对象类型。与事件表、条目表上的 object_type 取值一致。
const (
	ObjectArticle      = "article"
	ObjectTalk         = "talk"
	ObjectComment      = "comment"
	ObjectUser         = "user"
	ObjectAnnouncement = "announcement"
)

// 账号安全动作。
const (
	SecurityPasswordChanged = "password-changed"
	SecurityEmailChanged    = "email-changed"
)

// Payload 是通知事件的渲染载荷：业务事件存渲染所需的业务快照，公告事件只存内容引用。
//
// 字段按场景择用，未用到的留空即可；渲染时缺失的字段退化为通用文案，不阻断通知。
type Payload struct {
	// ActorName 触发者展示名。注册用户取昵称，访客兜底为「访客」，避免渲染出空串。
	ActorName string `json:"actor_name,omitempty"`
	// CommentExcerpt 评论摘要
	CommentExcerpt string `json:"comment_excerpt,omitempty"`
	// TopicType / TopicID 通知指向的主题，用于生成跳转地址
	TopicType string `json:"topic_type,omitempty"`
	TopicID   int64  `json:"topic_id,omitempty"`
	// ContentTitle 被评论、被点赞或被审核对象的标题
	ContentTitle string `json:"content_title,omitempty"`
	// Result 审核结果：rejected
	Result string `json:"result,omitempty"`
	// Reason 审核驳回原因
	Reason string `json:"reason,omitempty"`
	// Action 账号安全动作，见 Security* 常量
	Action string `json:"action,omitempty"`
	// AnnouncementID 公告内容引用。公告正文以内容表为准，不在载荷里复制。
	AnnouncementID int64 `json:"announcement_id,omitempty"`
}

// Content 是渲染结果：站内信、邮件、短信共用同一套文本，渠道只额外套外壳。
type Content struct {
	Title string
	Body  string
	Link  string
}
