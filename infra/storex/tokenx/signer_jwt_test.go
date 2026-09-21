package tokenx_test

import (
	"context"
	"testing"
	"time"

	"github.com/ve-weiyi/blog-cloud/infra/storex/tokenx"
)

// JWT 的 claims 里 iat / exp 都按 Unix 秒写入，HS256 又是确定性签名，
// 所以「同一秒内、同一份 claims」签出来的串完全相同。
//
// 现实后果：同一秒内调 Refresh 并不构成轮换 —— 新旧 access token 是同一个串，
// 因此旧串照样能通过校验。这不是缺陷，但很容易让人把断言建在错误的预期上。
func TestJWTSigner_Sign_IsDeterministicWithinSameSecond(t *testing.T) {
	signer := tokenx.NewJWTSigner([]byte("my-secret-key"), "test-issuer")
	ctx := context.Background()

	// 固定到秒，模拟「同一秒内的两次签发」
	now := time.Unix(1700000000, 0)
	claims := tokenx.TokenClaims{
		UserID:    "user1",
		DeviceID:  "device1",
		TokenType: "access",
		IssuedAt:  now,
		ExpiresAt: now.Add(2 * time.Hour),
	}

	first, err := signer.Sign(ctx, claims)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	second, err := signer.Sign(ctx, claims)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	if first != second {
		t.Error("同一秒内相同 claims 签出了不同的串 —— 若这是有意改动（如加入 jti 随机数），请一并调整依赖轮换的断言")
	}
}

func TestJWTSigner_Sign_ReturnsJWT(t *testing.T) {
	signer := tokenx.NewJWTSigner([]byte("my-secret-key"), "test-issuer")
	claims := tokenx.TokenClaims{
		UserID:    "user1",
		DeviceID:  "device1",
		IssuedAt:  time.Unix(1000, 0),
		ExpiresAt: time.Unix(1900, 0),
	}

	token, err := signer.Sign(context.Background(), claims)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
	// JWT has 3 dot-separated segments
	parts := 0
	for _, c := range token {
		if c == '.' {
			parts++
		}
	}
	if parts != 2 {
		t.Fatalf("expected JWT with 3 segments, got: %s", token)
	}
}

func TestJWTSigner_Verify_ReturnsClaims(t *testing.T) {
	signer := tokenx.NewJWTSigner([]byte("secret"), "iss")
	claims := tokenx.TokenClaims{
		UserID:    "user1",
		DeviceID:  "dev1",
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(time.Hour),
	}

	token, err := signer.Sign(context.Background(), claims)
	if err != nil {
		t.Fatal(err)
	}

	parsed, err := signer.Verify(context.Background(), token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parsed == nil {
		t.Fatal("expected non-nil claims for JWT")
	}
	if parsed.UserID != "user1" {
		t.Fatalf("expected UserID user1, got %s", parsed.UserID)
	}
	if parsed.DeviceID != "dev1" {
		t.Fatalf("expected DeviceID dev1, got %s", parsed.DeviceID)
	}
}

func TestJWTSigner_Verify_ExpiredToken(t *testing.T) {
	signer := tokenx.NewJWTSigner([]byte("secret"), "iss")
	claims := tokenx.TokenClaims{
		UserID: "u", DeviceID: "d",
		IssuedAt:  time.Now().Add(-2 * time.Hour),
		ExpiresAt: time.Now().Add(-1 * time.Hour), // expired
	}

	token, err := signer.Sign(context.Background(), claims)
	if err != nil {
		t.Fatal(err)
	}

	_, err = signer.Verify(context.Background(), token)
	if err != tokenx.ErrTokenExpired {
		t.Fatalf("expected ErrTokenExpired, got %v", err)
	}
}

func TestJWTSigner_Verify_InvalidToken(t *testing.T) {
	signer := tokenx.NewJWTSigner([]byte("secret"), "iss")

	_, err := signer.Verify(context.Background(), "not.a.valid.jwt")
	if err != tokenx.ErrTokenInvalid {
		t.Fatalf("expected ErrTokenInvalid, got %v", err)
	}
}

func TestJWTSigner_Verify_WrongKey(t *testing.T) {
	signer1 := tokenx.NewJWTSigner([]byte("secret1"), "iss")
	signer2 := tokenx.NewJWTSigner([]byte("secret2"), "iss")
	claims := tokenx.TokenClaims{
		UserID: "u", DeviceID: "d",
		IssuedAt: time.Now(), ExpiresAt: time.Now().Add(time.Hour),
	}

	token, _ := signer1.Sign(context.Background(), claims)
	_, err := signer2.Verify(context.Background(), token)
	if err != tokenx.ErrTokenInvalid {
		t.Fatalf("expected ErrTokenInvalid for wrong key, got %v", err)
	}
}
