package tokenx_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	tokenx2 "github.com/ve-weiyi/blog-cloud/infra/tokenx"
)

func TestIntegration_HMAC_MultiPoint_FullLifecycle(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	m, err := tokenx2.New(tokenx2.Config{
		Signer:          tokenx2.NewHMACSigner("integration-secret"),
		Store:           tokenx2.NewRedisStore(client),
		LoginMode:       tokenx2.MultiPoint,
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
	if err := m.Validate(ctx, "alice", "iphone", newPair.AccessToken); err != tokenx2.ErrTokenExpired {
		t.Fatalf("expected ErrTokenExpired after revoke, got %v", err)
	}
}

func TestIntegration_HMAC_SSO_SingleDeviceOnly(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	m, err := tokenx2.New(tokenx2.Config{
		Signer:    tokenx2.NewHMACSigner("sso-secret"),
		Store:     tokenx2.NewRedisStore(client),
		LoginMode: tokenx2.SinglePoint,
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
	if err := m.Validate(ctx, "bob", "laptop", pair1.AccessToken); err != tokenx2.ErrTokenExpired {
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

func TestIntegration_JWT_FullLifecycle(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	m, err := tokenx2.New(tokenx2.Config{
		Signer:    tokenx2.NewJWTSigner([]byte("jwt-secret-256bit-key!!"), "tokenx"),
		Store:     tokenx2.NewRedisStore(client),
		LoginMode: tokenx2.MultiPoint,
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
	if err := m.Validate(ctx, "evil", "browser", pair.AccessToken); err != tokenx2.ErrTokenInvalid {
		t.Fatalf("expected ErrTokenInvalid for wrong user, got %v", err)
	}

	// RevokeAll
	m.RevokeAll(ctx, "carol")
	if err := m.Validate(ctx, "carol", "browser", pair.AccessToken); err != tokenx2.ErrTokenExpired {
		t.Fatalf("expected ErrTokenExpired after RevokeAll, got %v", err)
	}
}
