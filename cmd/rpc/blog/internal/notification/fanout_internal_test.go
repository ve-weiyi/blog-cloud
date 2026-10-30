package notification

import (
	"testing"

	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
)

// 站内信没有投递地址：投递表只承载外部渠道，站内信的「投递」就是条目本身。
func TestDeliveryAddress_InboxHasNone(t *testing.T) {
	if got := deliveryAddress(notify.ChannelInbox, recipient{UserID: "u1"}); got != "" {
		t.Fatalf("站内信不应产生投递地址，得到 %q", got)
	}
	if got := deliveryAddress(notify.ChannelEmail, recipient{Email: "a@b.c"}); got != "a@b.c" {
		t.Fatalf("邮件地址 = %q，期望 a@b.c", got)
	}
	if got := deliveryAddress(notify.ChannelSms, recipient{Mobile: "13800000000"}); got != "13800000000" {
		t.Fatalf("短信地址 = %q", got)
	}
}
