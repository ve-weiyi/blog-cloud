package notify_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
)

const testSecret = "0123456789abcdef0123456789abcdef"

func newToken(t *testing.T, group notify.Group, expires time.Time) string {
	t.Helper()
	token, err := notify.SignUnsubscribe(testSecret, "u1", group, notify.ChannelEmail, expires)
	if err != nil {
		t.Fatalf("签发失败: %v", err)
	}
	return token
}

// 签发与校验闭环，令牌里带着收件人、分组与渠道。
func TestSignAndVerify(t *testing.T) {
	now := time.Now()
	token := newToken(t, notify.GroupComment, now.Add(time.Hour))

	claim, err := notify.VerifyUnsubscribe(testSecret, token, now)
	if err != nil {
		t.Fatalf("校验失败: %v", err)
	}
	if claim.UserID != "u1" || claim.Group != notify.GroupComment || claim.Channel != notify.ChannelEmail {
		t.Fatalf("载荷不符: %+v", claim)
	}
}

// 改动任意一段都应判签名不符。
func TestVerify_RejectsTamperedToken(t *testing.T) {
	now := time.Now()
	token := newToken(t, notify.GroupComment, now.Add(time.Hour))
	parts := strings.Split(token, ".")

	forged, err := notify.SignUnsubscribe(testSecret, "u2", notify.GroupComment, notify.ChannelEmail, now.Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	// 用别人的载荷配上原签名
	tampered := strings.Split(forged, ".")[0] + "." + parts[1]

	if _, err := notify.VerifyUnsubscribe(testSecret, tampered, now); !errors.Is(err, notify.ErrTokenInvalid) {
		t.Fatalf("篡改载荷应判签名不符，得到 %v", err)
	}
}

// 过期令牌不能再退订。
func TestVerify_RejectsExpiredToken(t *testing.T) {
	now := time.Now()
	token := newToken(t, notify.GroupComment, now.Add(-time.Minute))

	if _, err := notify.VerifyUnsubscribe(testSecret, token, now); !errors.Is(err, notify.ErrTokenExpired) {
		t.Fatalf("过期令牌应判过期，得到 %v", err)
	}
}

// 强制分组不签发令牌：邮件里根本不该出现它的退订链接。
func TestSign_RejectsMandatoryGroup(t *testing.T) {
	_, err := notify.SignUnsubscribe(testSecret, "u1", notify.GroupSecurity, notify.ChannelEmail, time.Now().Add(time.Hour))
	if err == nil {
		t.Fatal("强制分组不应签发退订令牌")
	}
}

// 换一把密钥签发的令牌同样校验不过。
func TestVerify_RejectsWrongSecret(t *testing.T) {
	now := time.Now()
	token := newToken(t, notify.GroupComment, now.Add(time.Hour))

	if _, err := notify.VerifyUnsubscribe(strings.Repeat("x", notify.MinSecretLen), token, now); !errors.Is(err, notify.ErrTokenInvalid) {
		t.Fatalf("换密钥应判签名不符，得到 %v", err)
	}
}

// 格式不对的令牌判为畸形，而不是签名不符。
func TestVerify_RejectsMalformedToken(t *testing.T) {
	if _, err := notify.VerifyUnsubscribe(testSecret, "not-a-token", time.Now()); !errors.Is(err, notify.ErrTokenMalformed) {
		t.Fatalf("格式错误应判畸形，得到 %v", err)
	}
}

// 没配密钥时两端都失败，不给「无签名也能退订」留后门。
func TestUnsubscribe_EmptySecretFailsClosed(t *testing.T) {
	if _, err := notify.SignUnsubscribe("", "u1", notify.GroupComment, notify.ChannelEmail, time.Now().Add(time.Hour)); err == nil {
		t.Fatal("空密钥不应签发成功")
	}
	token := newToken(t, notify.GroupComment, time.Now().Add(time.Hour))
	if _, err := notify.VerifyUnsubscribe("", token, time.Now()); err == nil {
		t.Fatal("空密钥不应校验通过")
	}
}
