// Package notify 定义通知的场景目录：哪些事件会产生通知、归属哪个偏好分组、
// 默认走哪些渠道、能否退订、在收件箱里怎么合并、保留多久，以及每个场景的文案。
package notify

import (
	"slices"
	"time"
)

// Level 通知等级，决定前端展示的强调程度。
type Level string

const (
	LevelInfo    Level = "info"
	LevelSuccess Level = "success"
	LevelWarning Level = "warning"
	LevelError   Level = "error"
)

// Channel 投递渠道。站内信是锚点渠道：它的条目同时是外部渠道的载体。
type Channel string

const (
	ChannelInbox Channel = "inbox"
	ChannelEmail Channel = "email"
	ChannelSms   Channel = "sms"
	ChannelPush  Channel = "push"
)

// Group 偏好分组：用户退订时面对的粒度，不随子场景膨胀。
type Group string

const (
	GroupAnnouncement Group = "announcement"
	GroupComment      Group = "comment"
	GroupLike         Group = "like"
	GroupReview       Group = "review"
	GroupSecurity     Group = "security"
)

// Grouping 展示分组：收件箱里如何合并条目。
type Grouping string

const (
	// GroupingNone 每条条目自成一组
	GroupingNone Grouping = "none"
	// GroupingObject 同一对象的条目合并为一组
	GroupingObject Grouping = "object"
)

// 场景标识。新增场景只需要在这里加一项常量与一条目录项，不改表结构。
const (
	SceneAnnouncement     = "announcement"
	SceneContentCommented = "content.commented"
	SceneCommentReplied   = "comment.replied"
	SceneCommentLiked     = "comment.liked"
	SceneArticleLiked     = "article.liked"
	SceneTalkLiked        = "talk.liked"
	SceneContentReviewed  = "content.reviewed"
	SceneAccountSecurity  = "account.security"
)

// Scene 一个场景的元数据。
type Scene struct {
	Key      string
	Group    Group
	Level    Level
	Channels []Channel
	// Mandatory 为真表示不可退订，偏好表对它无效
	Mandatory bool
	Grouping  Grouping
	// ReadTTL / UnreadTTL 保留期，已读先于未读清理
	ReadTTL   time.Duration
	UnreadTTL time.Duration
}

const day = 24 * time.Hour

// catalog 是场景目录。没有列在这里的场景标识在扇出时视为编程错误，直接失败。
var catalog = map[string]Scene{
	SceneAnnouncement: {
		Key:       SceneAnnouncement,
		Group:     GroupAnnouncement,
		Level:     LevelInfo,
		Channels:  []Channel{ChannelInbox},
		Mandatory: true,
		Grouping:  GroupingNone,
		ReadTTL:   180 * day,
		UnreadTTL: 365 * day,
	},
	SceneContentCommented: {
		Key:       SceneContentCommented,
		Group:     GroupComment,
		Level:     LevelInfo,
		Channels:  []Channel{ChannelInbox},
		Grouping:  GroupingObject,
		ReadTTL:   90 * day,
		UnreadTTL: 180 * day,
	},
	SceneCommentReplied: {
		Key:       SceneCommentReplied,
		Group:     GroupComment,
		Level:     LevelInfo,
		Channels:  []Channel{ChannelInbox},
		Grouping:  GroupingNone,
		ReadTTL:   90 * day,
		UnreadTTL: 180 * day,
	},
	SceneCommentLiked: {
		Key:       SceneCommentLiked,
		Group:     GroupLike,
		Level:     LevelInfo,
		Channels:  []Channel{ChannelInbox},
		Grouping:  GroupingObject,
		ReadTTL:   90 * day,
		UnreadTTL: 180 * day,
	},
	SceneArticleLiked: {
		Key:       SceneArticleLiked,
		Group:     GroupLike,
		Level:     LevelInfo,
		Channels:  []Channel{ChannelInbox},
		Grouping:  GroupingObject,
		ReadTTL:   90 * day,
		UnreadTTL: 180 * day,
	},
	SceneTalkLiked: {
		Key:       SceneTalkLiked,
		Group:     GroupLike,
		Level:     LevelInfo,
		Channels:  []Channel{ChannelInbox},
		Grouping:  GroupingObject,
		ReadTTL:   90 * day,
		UnreadTTL: 180 * day,
	},
	SceneContentReviewed: {
		Key:       SceneContentReviewed,
		Group:     GroupReview,
		Level:     LevelWarning,
		Channels:  []Channel{ChannelInbox, ChannelEmail},
		Grouping:  GroupingNone,
		ReadTTL:   90 * day,
		UnreadTTL: 180 * day,
	},
	SceneAccountSecurity: {
		Key:       SceneAccountSecurity,
		Group:     GroupSecurity,
		Level:     LevelWarning,
		Channels:  []Channel{ChannelInbox, ChannelEmail},
		Mandatory: true,
		Grouping:  GroupingNone,
		ReadTTL:   180 * day,
		UnreadTTL: 365 * day,
	},
}

// Lookup 返回场景定义；未注册的场景返回 false，调用方应视为编程错误直接失败。
func Lookup(key string) (Scene, bool) {
	s, ok := catalog[key]
	return s, ok
}

// All 返回全部场景，顺序不稳定，调用方需要稳定顺序时自行排序。
func All() []Scene {
	out := make([]Scene, 0, len(catalog))
	for _, s := range catalog {
		out = append(out, s)
	}
	return out
}

// Groups 返回全部偏好分组。
func Groups() []Group {
	return []Group{GroupAnnouncement, GroupComment, GroupLike, GroupReview, GroupSecurity}
}

// ScenesOfGroups 把偏好分组展开成场景集合，供收件箱按分组过滤。
// 空集合表示全部场景；结果排序后返回，让同一请求生成的 SQL 参数稳定。
func ScenesOfGroups(groups []Group) []string {
	want := make(map[Group]struct{}, len(groups))
	for _, g := range groups {
		want[g] = struct{}{}
	}

	out := make([]string, 0, len(catalog))
	for k, s := range catalog {
		if len(want) > 0 {
			if _, ok := want[s.Group]; !ok {
				continue
			}
		}
		out = append(out, k)
	}
	slices.Sort(out)
	return out
}
