package notify_test

import (
	"testing"

	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
)

// 每个场景常量都必须有目录项，否则扇出时会当成编程错误。
func TestLookup_AllSceneConstantsRegistered(t *testing.T) {
	scenes := []string{
		notify.SceneAnnouncement,
		notify.SceneContentCommented,
		notify.SceneCommentReplied,
		notify.SceneCommentLiked,
		notify.SceneArticleLiked,
		notify.SceneTalkLiked,
		notify.SceneContentReviewed,
		notify.SceneAccountSecurity,
	}

	for _, key := range scenes {
		s, ok := notify.Lookup(key)
		if !ok {
			t.Fatalf("场景 %q 未注册到目录", key)
		}
		if s.Key != key {
			t.Fatalf("场景 %q 的 Key 是 %q，两者必须一致", key, s.Key)
		}
	}
}

// 目录项的取值必须落在已知枚举内，且保留期满足「未读不短于已读」。
func TestCatalog_WellFormed(t *testing.T) {
	validGroups := map[notify.Group]bool{}
	for _, g := range notify.Groups() {
		validGroups[g] = true
	}

	for _, s := range notify.All() {
		if !validGroups[s.Group] {
			t.Errorf("场景 %q 的偏好分组 %q 不在 Groups() 内", s.Key, s.Group)
		}
		if len(s.Channels) == 0 {
			t.Errorf("场景 %q 没有默认渠道", s.Key)
		}
		if s.Grouping != notify.GroupingNone && s.Grouping != notify.GroupingObject {
			t.Errorf("场景 %q 的展示分组 %q 非法", s.Key, s.Grouping)
		}
		if s.UnreadTTL < s.ReadTTL {
			t.Errorf("场景 %q 的未读保留期短于已读", s.Key)
		}
	}
}

// 强制分组不可退订：公告与账号安全是底线。
func TestCatalog_MandatoryScenes(t *testing.T) {
	for _, key := range []string{notify.SceneAnnouncement, notify.SceneAccountSecurity} {
		s, _ := notify.Lookup(key)
		if !s.Mandatory {
			t.Errorf("场景 %q 必须标记为 Mandatory", key)
		}
	}
}

// 偏好分组展开成场景集合；空集合表示全部。
func TestScenesOfGroups(t *testing.T) {
	all := notify.ScenesOfGroups(nil)
	if len(all) != len(notify.All()) {
		t.Fatalf("空分组应展开为全部场景，得到 %d 个，期望 %d 个", len(all), len(notify.All()))
	}

	likes := notify.ScenesOfGroups([]notify.Group{notify.GroupLike})
	want := map[string]bool{
		notify.SceneCommentLiked: true,
		notify.SceneArticleLiked: true,
		notify.SceneTalkLiked:    true,
	}
	if len(likes) != len(want) {
		t.Fatalf("点赞分组应展开为 %d 个场景，得到 %v", len(want), likes)
	}
	for _, key := range likes {
		if !want[key] {
			t.Errorf("点赞分组里出现意料之外的场景 %q", key)
		}
	}
}

// 业务场景按触发者人数切换单人与多人文案。
func TestRender_ActorCount(t *testing.T) {
	p := notify.Payload{ActorName: "小明", ContentTitle: "Go 并发", TopicType: notify.ObjectArticle, TopicID: 7}

	single, ok := notify.Render(notify.SceneArticleLiked, p, 1, "zh")
	if !ok {
		t.Fatal("文章被点赞应由本包渲染")
	}
	if want := "小明 赞了你的文章《Go 并发》"; single.Title != want {
		t.Errorf("单人标题 = %q，期望 %q", single.Title, want)
	}
	if single.Link != "/article/7" {
		t.Errorf("跳转 = %q，期望 /article/7", single.Link)
	}

	multi, _ := notify.Render(notify.SceneArticleLiked, p, 5, "zh")
	if want := "小明 等 5 人 赞了你的文章《Go 并发》"; multi.Title != want {
		t.Errorf("多人标题 = %q，期望 %q", multi.Title, want)
	}
}

// 公告不在本包渲染：正文以内容表为准，调用方负责。
func TestRender_AnnouncementDelegated(t *testing.T) {
	if _, ok := notify.Render(notify.SceneAnnouncement, notify.Payload{}, 1, "zh"); ok {
		t.Fatal("公告应交给调用方渲染")
	}
}

// 场景已从目录移除、历史条目仍在时，不能渲染失败。
func TestRender_UnknownSceneFallsBack(t *testing.T) {
	c, ok := notify.Render("legacy.scene", notify.Payload{ActorName: "小明"}, 1, "zh")
	if !ok || c.Title == "" {
		t.Fatalf("未知场景应回退到通用文案，得到 %+v ok=%v", c, ok)
	}
}

// 缺失字段不能让文案出现空书名号或空触发者。
func TestRender_MissingFieldsDegrade(t *testing.T) {
	c, _ := notify.Render(notify.SceneArticleLiked, notify.Payload{}, 1, "zh")
	if c.Title != "有人 赞了你的文章" {
		t.Errorf("缺触发者与标题时标题 = %q", c.Title)
	}
	if c.Link != "" {
		t.Errorf("缺主题信息时不应生成跳转，得到 %q", c.Link)
	}
}
