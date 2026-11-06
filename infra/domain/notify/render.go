package notify

import "fmt"

// Render 按场景与载荷渲染标题、正文与跳转。
//
// actorCount 是展示组内的人数：1 表示单条，大于 1 时文案改为「A 等 N 人」。
// 返回值 ok 为 false 表示该场景不在本包渲染 —— 目前只有公告，它的正文以内容表为准。
//
// lang 当前只有一种语言，保留参数以免日后加语言时改动一圈调用方。
func Render(scene string, p Payload, actorCount int64, lang string) (Content, bool) {
	if scene == SceneAnnouncement {
		return Content{}, false
	}

	who := actorPhrase(p.ActorName, actorCount)
	link := topicLink(p.TopicType, p.TopicID)

	switch scene {
	case SceneContentCommented:
		return Content{
			Title: "你的" + topicNoun(p.TopicType) + "有新评论",
			Body:  fmt.Sprintf("%s 评论了%s：%s", who, quote(p.ContentTitle), p.CommentExcerpt),
			Link:  link,
		}, true
	case SceneCommentReplied:
		return Content{
			Title: fmt.Sprintf("%s 回复了你的评论", who),
			Body:  p.CommentExcerpt,
			Link:  link,
		}, true
	case SceneCommentLiked:
		return Content{
			Title: fmt.Sprintf("%s 赞了你的评论", who),
			Body:  p.CommentExcerpt,
			Link:  link,
		}, true
	case SceneArticleLiked:
		return Content{
			Title: fmt.Sprintf("%s 赞了你的文章%s", who, quote(p.ContentTitle)),
			Link:  link,
		}, true
	case SceneTalkLiked:
		return Content{
			Title: fmt.Sprintf("%s 赞了你的说说", who),
			Link:  link,
		}, true
	case SceneContentReviewed:
		body := p.Reason
		if body == "" {
			body = "请修改后重新提交"
		}
		return Content{
			Title: fmt.Sprintf("你的%s未通过审核", topicNoun(p.TopicType)),
			Body:  body,
			Link:  link,
		}, true
	case SceneAccountSecurity:
		return Content{
			Title: securityTitle(p.Action),
			Body:  "如果这不是你本人的操作，请立即修改密码。",
		}, true
	default:
		// 场景已从目录移除但历史条目仍在：给一个通用文案，不让整页渲染失败
		return Content{Title: "你有一条新通知", Body: p.CommentExcerpt, Link: link}, true
	}
}

// actorPhrase 生成触发者称谓；多人为「A 等 N 人」，缺名时退化为「有人」。
func actorPhrase(name string, count int64) string {
	if name == "" {
		name = "有人"
	}
	if count > 1 {
		return fmt.Sprintf("%s 等 %d 人", name, count)
	}
	return name
}

// topicNoun 把对象类型换成文案里的名词。
func topicNoun(topicType string) string {
	switch topicType {
	case ObjectArticle:
		return "文章"
	case ObjectTalk:
		return "说说"
	default:
		return "内容"
	}
}

// quote 给标题加书名号；标题为空时返回空串，避免渲染出「《》」。
func quote(title string) string {
	if title == "" {
		return ""
	}
	return "《" + title + "》"
}

// topicLink 生成跳转地址。主题信息缺失时返回空串，客户端停留在当前页。
func topicLink(topicType string, topicID int64) string {
	if topicID <= 0 {
		return ""
	}
	switch topicType {
	case ObjectArticle:
		return fmt.Sprintf("/article/%d", topicID)
	case ObjectTalk:
		return fmt.Sprintf("/talk/%d", topicID)
	default:
		return ""
	}
}

// securityTitle 把安全动作换成标题。
func securityTitle(action string) string {
	switch action {
	case SecurityPasswordChanged:
		return "账号密码已修改"
	case SecurityEmailChanged:
		return "账号邮箱已变更"
	default:
		return "账号安全提醒"
	}
}
