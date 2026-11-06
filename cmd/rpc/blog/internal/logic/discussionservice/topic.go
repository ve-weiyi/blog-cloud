package discussionservicelogic

import (
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
	"github.com/ve-weiyi/blog-cloud/model"
)

// topicOf 把评论类型映射成通知里的主题类型与 ID。
// 友链评论没有可跳转的主题，返回空值，由渲染层退化为不生成跳转。
func topicOf(c *model.TComment) (string, int64) {
	switch c.Type {
	case enums.CommentTypeArticle:
		return notify.ObjectArticle, c.TopicId
	case enums.CommentTypeTalk:
		return notify.ObjectTalk, c.TopicId
	default:
		return "", 0
	}
}

// excerpt 截断正文，避免把整段内容塞进通知载荷。
func excerpt(s string, limit int) string {
	r := []rune(s)
	if len(r) <= limit {
		return s
	}
	return string(r[:limit]) + "…"
}
