package captchax

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	"github.com/ve-weiyi/blog-cloud/infra/storex"
)

func newTestStore(t *testing.T, prefix string) (*Store, *miniredis.Miniredis) {
	t.Helper()

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("start miniredis: %v", err)
	}
	t.Cleanup(mr.Close)

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = rdb.Close() })

	return New(storex.NewRedisStore(rdb), prefix, 15*time.Minute), mr
}

// TestVerifyCaptcha 覆盖一次性验证码的完整生命周期。
func TestVerifyCaptcha(t *testing.T) {
	store, mr := newTestStore(t, "blog:admin:captcha:")
	ctx := context.Background()

	key, _, answer, err := store.GetMathImageCaptcha(0, 0)
	if err != nil {
		t.Fatalf("generate captcha: %v", err)
	}

	// 答案落在带前缀的 key 上
	if got, _ := mr.Get("blog:admin:captcha:" + key); got != answer {
		t.Fatalf("stored answer = %q, want %q", got, answer)
	}

	// 错误答案不通过，且 key 被消费
	if ok, err := store.VerifyCaptcha(ctx, key, answer+"x"); err != nil || ok {
		t.Fatalf("wrong answer: ok=%v err=%v, want false/nil", ok, err)
	}
	if mr.Exists("blog:admin:captcha:" + key) {
		t.Fatal("captcha key should be deleted after verification")
	}

	// 正确答案通过
	key, _, answer, err = store.GetMathImageCaptcha(0, 0)
	if err != nil {
		t.Fatalf("generate captcha: %v", err)
	}
	if ok, err := store.VerifyCaptcha(ctx, key, answer); err != nil || !ok {
		t.Fatalf("correct answer: ok=%v err=%v, want true/nil", ok, err)
	}

	// 一次性：同一 key 不可复用
	if ok, _ := store.VerifyCaptcha(ctx, key, answer); ok {
		t.Fatal("captcha key should not be reusable")
	}
}

// TestVerifyCaptchaEmpty 空参数直接判否，不触碰存储。
func TestVerifyCaptchaEmpty(t *testing.T) {
	store, _ := newTestStore(t, "blog:admin:captcha:")

	if ok, err := store.VerifyCaptcha(context.Background(), "", ""); ok || err != nil {
		t.Fatalf("empty input: ok=%v err=%v, want false/nil", ok, err)
	}
}

// TestVerifyCaptchaPrefixIsolation 不同服务前缀互不可见。
func TestVerifyCaptchaPrefixIsolation(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("start miniredis: %v", err)
	}
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer func() { _ = rdb.Close() }()

	admin := New(storex.NewRedisStore(rdb), "blog:admin:captcha:", 15*time.Minute)
	app := New(storex.NewRedisStore(rdb), "blog:app:captcha:", 15*time.Minute)
	ctx := context.Background()

	key, _, answer, err := admin.GetMathImageCaptcha(0, 0)
	if err != nil {
		t.Fatalf("generate captcha: %v", err)
	}

	if ok, _ := app.VerifyCaptcha(ctx, key, answer); ok {
		t.Fatal("answer issued under admin prefix should not verify under app prefix")
	}
}
