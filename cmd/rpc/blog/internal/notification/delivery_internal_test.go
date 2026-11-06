package notification

import (
	"strings"
	"testing"

	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
)

// 正文里可能夹带用户内容，进邮件前必须转义，避免注入到 HTML。
func TestEmailShell_EscapesContent(t *testing.T) {
	got := emailShell("<b>标题</b>", `内容 <script>alert(1)</script>`, "")

	if strings.Contains(got, "<script>") {
		t.Fatalf("脚本未转义: %s", got)
	}
	if !strings.Contains(got, "&lt;script&gt;") {
		t.Fatalf("应转义为实体: %s", got)
	}
	if !strings.Contains(got, "&lt;b&gt;标题&lt;/b&gt;") {
		t.Fatalf("标题应转义: %s", got)
	}
}

// 退订链接只在「可退订分组 + 密钥/站点地址都配好」时才渲染。
func TestUnsubscribeLink(t *testing.T) {
	commentScene, _ := notify.Lookup(notify.SceneContentReviewed)
	securityScene, _ := notify.Lookup(notify.SceneAccountSecurity)

	secret := strings.Repeat("k", notify.MinSecretLen)
	configured := &DeliveryWorker{svc: &Service{unsubscribeSecret: secret, siteURL: "https://blog.test/"}}
	if got := configured.unsubscribeLink("u1", commentScene, notify.ChannelEmail); !strings.Contains(got, "https://blog.test/unsubscribe?token=") {
		t.Fatalf("已配置时应生成链接，得到 %q", got)
	}

	// 强制分组不签发
	if got := configured.unsubscribeLink("u1", securityScene, notify.ChannelEmail); got != "" {
		t.Fatalf("强制分组不应生成退订链接，得到 %q", got)
	}

	// 缺配置时宁可不渲染
	unconfigured := &DeliveryWorker{svc: &Service{}}
	if got := unconfigured.unsubscribeLink("u1", commentScene, notify.ChannelEmail); got != "" {
		t.Fatalf("缺配置时不应生成链接，得到 %q", got)
	}

	// 占位值或短口令按未配置处理：猜得出的密钥等于没有密钥
	short := &DeliveryWorker{svc: &Service{unsubscribeSecret: "CHANGE_ME", siteURL: "https://blog.test/"}}
	if got := short.unsubscribeLink("u1", commentScene, notify.ChannelEmail); got != "" {
		t.Fatalf("密钥过短时不应生成链接，得到 %q", got)
	}
}

// 退订链接进邮件时同样要转义，避免链接里的特殊字符破坏 HTML。
func TestEmailShell_RendersUnsubscribeFooter(t *testing.T) {
	got := emailShell("标题", "正文", "https://blog.test/unsubscribe?token=a&b")
	if !strings.Contains(got, "点此退订") {
		t.Fatal("应渲染退订页脚")
	}
	if !strings.Contains(got, "a&amp;b") {
		t.Fatalf("链接应转义: %s", got)
	}
}
