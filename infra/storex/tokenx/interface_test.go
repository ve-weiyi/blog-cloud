package tokenx_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	"github.com/ve-weiyi/blog-cloud/infra/storex"
	"github.com/ve-weiyi/blog-cloud/infra/storex/tokenx"
)

// newRedisStores 起一个 miniredis，返回接在它上面的 KV 与 member 存储。
func newRedisStores(t *testing.T) (storex.KVStore, storex.MemberStore) {
	t.Helper()

	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	return storex.NewRedisStore(client), storex.NewRedisMemberStore(client)
}

func TestIntegration_HMAC_MultiPoint_FullLifecycle(t *testing.T) {
	kv, devices := newRedisStores(t)

	m, err := tokenx.New(tokenx.Config{
		Signer:          tokenx.NewHMACSigner("integration-secret"),
		Store:           kv,
		Devices:         devices,
		LoginMode:       tokenx.MultiPoint,
		KeyPrefix:       "myapp",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	// Generate
	pair, err := m.Generate(ctx, "alice", "iphone")
	if err != nil {
		t.Fatal(err)
	}

	// Validate
	if err := m.Validate(ctx, "alice", "iphone", pair.AccessToken); err != nil {
		t.Fatal(err)
	}

	// Refresh
	newPair, err := m.Refresh(ctx, "alice", "iphone", pair.RefreshToken)
	if err != nil {
		t.Fatal(err)
	}

	// New access token should validate
	if err := m.Validate(ctx, "alice", "iphone", newPair.AccessToken); err != nil {
		t.Fatal(err)
	}

	// List sessions
	sessions, err := m.ListSessions(ctx, "alice")
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}

	// Revoke
	if err := m.Revoke(ctx, "alice", "iphone"); err != nil {
		t.Fatal(err)
	}

	// Validate after revoke
	if err := m.Validate(ctx, "alice", "iphone", newPair.AccessToken); err != tokenx.ErrTokenExpired {
		t.Fatalf("expected ErrTokenExpired after revoke, got %v", err)
	}
}

func TestIntegration_HMAC_SSO_SingleDeviceOnly(t *testing.T) {
	kv, devices := newRedisStores(t)

	m, err := tokenx.New(tokenx.Config{
		Signer:    tokenx.NewHMACSigner("sso-secret"),
		Store:     kv,
		Devices:   devices,
		LoginMode: tokenx.SinglePoint,
		KeyPrefix: "sso",
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	// Login on device1
	pair1, _ := m.Generate(ctx, "bob", "laptop")
	// Login on device2 → auto-revoke device1
	pair2, _ := m.Generate(ctx, "bob", "tablet")

	// device1 should be revoked
	if err := m.Validate(ctx, "bob", "laptop", pair1.AccessToken); err != tokenx.ErrTokenExpired {
		t.Fatalf("expected device1 revoked via SSO, got %v", err)
	}

	// device2 should be active
	if err := m.Validate(ctx, "bob", "tablet", pair2.AccessToken); err != nil {
		t.Fatalf("expected device2 valid, got %v", err)
	}

	// Only 1 session
	sessions, _ := m.ListSessions(ctx, "bob")
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session under SSO, got %d", len(sessions))
	}
}

// 这条走的是生产实际的组合：api/app 与 api/admin 都用 JWTSigner + SinglePoint。
//
// 单点登录的吊销完全依赖「store 里查不到这个 key」，而这次重构正好把 Get 的契约
// 改成了 (value, found, err)。所以这条路径必须单独钉住，不能只靠 HMAC 的用例覆盖。
func TestIntegration_JWT_SSO_MirrorsProductionConfig(t *testing.T) {
	kv, devices := newRedisStores(t)

	m, err := tokenx.New(tokenx.Config{
		Signer:          tokenx.NewJWTSigner([]byte("prod-shaped-secret"), "blog-cloud"),
		Store:           kv,
		Devices:         devices,
		LoginMode:       tokenx.SinglePoint,
		KeyPrefix:       "blog:app:token:",
		AccessTokenTTL:  2 * time.Hour,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	// 首次登录 → 校验通过
	first, err := m.Generate(ctx, "alice", "iphone")
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Validate(ctx, "alice", "iphone", first.AccessToken); err != nil {
		t.Fatalf("首次登录后校验失败：%v", err)
	}

	// 换设备登录 → 单点把上一台踢掉
	second, err := m.Generate(ctx, "alice", "ipad")
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Validate(ctx, "alice", "iphone", first.AccessToken); err != tokenx.ErrTokenExpired {
		t.Errorf("上一台设备未被踢掉：%v", err)
	}
	if err := m.Validate(ctx, "alice", "ipad", second.AccessToken); err != nil {
		t.Fatalf("新设备校验失败：%v", err)
	}

	// 续期 → 新令牌可校验。
	//
	// 这里不断言「旧 access token 立即失效」：JWT 的 iat/exp 只到秒级，
	// 同一秒内同一份 claims 签出来是同一个串，旧串就是当前存着的那串。
	// 轮换使旧令牌失效这条语义由 TestManager_Refresh_OldTokenInvalidated 覆盖
	// （HMAC + 内存 store，两次签名结果必然不同）。
	refreshed, err := m.Refresh(ctx, "alice", "ipad", second.RefreshToken)
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Validate(ctx, "alice", "ipad", refreshed.AccessToken); err != nil {
		t.Fatalf("续期后校验失败：%v", err)
	}

	// 单点模式下会话清单只应有当前这一台
	sessions, err := m.ListSessions(ctx, "alice")
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 {
		t.Fatalf("会话数 = %d, want 1：%+v", len(sessions), sessions)
	}
	if sessions[0].DeviceID != "ipad" {
		t.Errorf("剩余会话的 DeviceID = %q, want %q", sessions[0].DeviceID, "ipad")
	}

	// 登出 → 校验必须失败
	if err := m.Revoke(ctx, "alice", "ipad"); err != nil {
		t.Fatal(err)
	}
	if err := m.Validate(ctx, "alice", "ipad", refreshed.AccessToken); err != tokenx.ErrTokenExpired {
		t.Errorf("登出后仍能通过校验：%v", err)
	}
}

func TestIntegration_JWT_FullLifecycle(t *testing.T) {
	kv, devices := newRedisStores(t)

	m, err := tokenx.New(tokenx.Config{
		Signer:    tokenx.NewJWTSigner([]byte("jwt-secret-256bit-key!!"), "tokenx"),
		Store:     kv,
		Devices:   devices,
		LoginMode: tokenx.MultiPoint,
		KeyPrefix: "jwtapp",
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	// Generate
	pair, err := m.Generate(ctx, "carol", "browser")
	if err != nil {
		t.Fatal(err)
	}

	// Validate with correct triple
	if err := m.Validate(ctx, "carol", "browser", pair.AccessToken); err != nil {
		t.Fatal(err)
	}

	// Validate with wrong user (JWT claims cross-check)
	if err := m.Validate(ctx, "evil", "browser", pair.AccessToken); err != tokenx.ErrTokenInvalid {
		t.Fatalf("expected ErrTokenInvalid for wrong user, got %v", err)
	}

	// RevokeAll
	m.RevokeAll(ctx, "carol")
	if err := m.Validate(ctx, "carol", "browser", pair.AccessToken); err != tokenx.ErrTokenExpired {
		t.Fatalf("expected ErrTokenExpired after RevokeAll, got %v", err)
	}
}
