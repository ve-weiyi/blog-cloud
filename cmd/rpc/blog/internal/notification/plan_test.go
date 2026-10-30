package notification_test

import (
	"strings"
	"testing"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/notification"
	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
)

func sceneOf(t *testing.T, key string) notify.Scene {
	t.Helper()
	s, ok := notify.Lookup(key)
	if !ok {
		t.Fatalf("场景 %q 未注册", key)
	}
	return s
}

// 可分组场景按对象归并，不分组场景一条一组。
func TestGroupKey(t *testing.T) {
	liked := sceneOf(t, notify.SceneArticleLiked)
	if got := notification.GroupKey(liked, notify.ObjectArticle, 7, 99); got != "article.liked:article:7" {
		t.Errorf("可分组场景组键 = %q", got)
	}

	replied := sceneOf(t, notify.SceneCommentReplied)
	if got := notification.GroupKey(replied, notify.ObjectComment, 7, 99); got != "comment.replied:e:99" {
		t.Errorf("不分组场景组键 = %q", got)
	}
}

// 可分组场景按触发者去重：同一个人反复互动落在同一行，组大小才等于人数。
func TestDedupeKey(t *testing.T) {
	liked := sceneOf(t, notify.SceneArticleLiked)
	group := notification.GroupKey(liked, notify.ObjectArticle, 7, 99)

	if a, b := notification.DedupeKey(liked, group, "u1"), notification.DedupeKey(liked, group, "u1"); a != b {
		t.Errorf("同一触发者应得到同一去重键：%q vs %q", a, b)
	}
	if a, b := notification.DedupeKey(liked, group, "u1"), notification.DedupeKey(liked, group, "u2"); a == b {
		t.Errorf("不同触发者不应共用去重键：%q", a)
	}

	replied := sceneOf(t, notify.SceneCommentReplied)
	group = notification.GroupKey(replied, notify.ObjectComment, 7, 99)
	if got := notification.DedupeKey(replied, group, "u1"); got != group {
		t.Errorf("不分组场景的去重键应与组键相同，得到 %q，期望 %q", got, group)
	}
}

// 键长必须落在 dedupe_key 列宽内，否则插入会被截断或报错。
func TestKeys_FitColumnWidth(t *testing.T) {
	worst := notify.Scene{
		Key:      strings.Repeat("s", 48),
		Grouping: notify.GroupingObject,
	}
	group := notification.GroupKey(worst, strings.Repeat("o", 32), 9223372036854775807, 0)
	dedupe := notification.DedupeKey(worst, group, strings.Repeat("u", 64))

	if n := len(dedupe); n > 191 {
		t.Fatalf("去重键长 %d，超出列宽 191：%q", n, dedupe)
	}
}

// 偏好：强制分组不看偏好；站内信是锚点，关掉它就什么都不发。
func TestChannels(t *testing.T) {
	security := sceneOf(t, notify.SceneAccountSecurity)
	if got := notification.Channels(security, map[notify.Channel]bool{notify.ChannelInbox: false}); len(got) != len(security.Channels) {
		t.Errorf("强制分组应忽略偏好，得到 %v", got)
	}

	reviewed := sceneOf(t, notify.SceneContentReviewed)

	if got := notification.Channels(reviewed, nil); len(got) != 2 {
		t.Errorf("默认应走站内信与邮件，得到 %v", got)
	}
	if got := notification.Channels(reviewed, map[notify.Channel]bool{notify.ChannelInbox: false}); len(got) != 0 {
		t.Errorf("关闭站内信应关闭全部渠道，得到 %v", got)
	}
	if got := notification.Channels(reviewed, map[notify.Channel]bool{notify.ChannelEmail: false}); len(got) != 1 || got[0] != notify.ChannelInbox {
		t.Errorf("只关邮件时应保留站内信，得到 %v", got)
	}
}

// 受众清洗：去空白、去空值、去重、剔除触发者，并保持首次出现的顺序。
func TestNormalizeRecipients(t *testing.T) {
	got := notification.NormalizeRecipients([]string{" u1 ", "", "u2", "u1", "actor", "  "}, "actor")
	want := []string{"u1", "u2"}

	if len(got) != len(want) {
		t.Fatalf("收件人 = %v，期望 %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("收件人 = %v，期望 %v", got, want)
		}
	}
}

// 受众类型为空必须判非法，不能被当成全站广播。
func TestValidAudience(t *testing.T) {
	if notification.ValidAudience("") {
		t.Error("空值应判非法")
	}
	if !notification.ValidAudience(notification.AudienceAll) {
		t.Error("all 应合法")
	}
	if !notification.ValidAudience(notification.AudienceUserIDs) {
		t.Error("user_ids 应合法")
	}
}

// 条目计划：偏好全关的收件人不产生条目；可分组场景同组不同人共用组键、去重键按人区分。
func TestPlanEntries(t *testing.T) {
	liked := sceneOf(t, notify.SceneArticleLiked)
	prefs := map[string]map[notify.Channel]bool{
		"u2": {notify.ChannelInbox: false},
	}

	entries := notification.PlanEntries(liked, notify.ObjectArticle, 7, 99, "actor",
		[]string{"u1", "u2", "u3"}, prefs)

	if len(entries) != 2 {
		t.Fatalf("关闭站内信的收件人应被跳过，得到 %d 条：%+v", len(entries), entries)
	}
	if entries[0].GroupKey != entries[1].GroupKey {
		t.Errorf("同一对象的组键应相同：%q vs %q", entries[0].GroupKey, entries[1].GroupKey)
	}
	// 去重键描述的是「谁对哪个对象做了什么」，同一触发者下同一事件的所有收件人共用；
	// 唯一性由「收件人 + 去重键」两列共同保证。
	if entries[0].DedupeKey != entries[1].DedupeKey {
		t.Errorf("同一触发者的去重键应相同：%q vs %q", entries[0].DedupeKey, entries[1].DedupeKey)
	}
	if entries[0].RecipientID == entries[1].RecipientID {
		t.Errorf("两条条目的收件人不应相同：%q", entries[0].RecipientID)
	}
}

// 撤回互动必须命中扇出写下的那一行：键的算法与写入路径一致，且与事件 id 无关。
func TestInteractionKey(t *testing.T) {
	commented := sceneOf(t, notify.SceneContentCommented)

	entries := notification.PlanEntries(commented, notify.ObjectArticle, 7, 99, "u1", []string{"author"}, nil)
	if len(entries) != 1 {
		t.Fatalf("应有 1 条计划，得到 %+v", entries)
	}
	if got := notification.InteractionKey(commented, notify.ObjectArticle, 7, "u1"); got != entries[0].DedupeKey {
		t.Errorf("撤回键 %q 与写入键 %q 不一致", got, entries[0].DedupeKey)
	}

	// 不分组场景没有「一次互动」的定位语义，撤回路径不处理它；这里只守住键本身稳定。
	replied := sceneOf(t, notify.SceneCommentReplied)
	if a, b := notification.InteractionKey(replied, notify.ObjectComment, 7, "u1"),
		notification.InteractionKey(replied, notify.ObjectComment, 7, "u1"); a != b {
		t.Errorf("同一互动的键应稳定：%q vs %q", a, b)
	}
}

// 只关邮件时条目保留，渠道里不再有邮件。
func TestPlanEntries_ChannelFiltering(t *testing.T) {
	reviewed := sceneOf(t, notify.SceneContentReviewed)
	prefs := map[string]map[notify.Channel]bool{
		"u1": {notify.ChannelEmail: false},
	}

	entries := notification.PlanEntries(reviewed, notify.ObjectComment, 7, 99, "", []string{"u1"}, prefs)
	if len(entries) != 1 {
		t.Fatalf("应保留 1 条，得到 %+v", entries)
	}
	if len(entries[0].External) != 0 {
		t.Errorf("站内信不应进入外部渠道，得到 %v", entries[0].External)
	}
}
