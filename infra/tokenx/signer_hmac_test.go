package tokenx_test

import (
	"context"
	"testing"
	"time"

	tokenx2 "github.com/ve-weiyi/blog-cloud/infra/tokenx"
)

func TestHMACSigner_Sign_ReturnsHexToken(t *testing.T) {
	signer := tokenx2.NewHMACSigner("my-secret")
	claims := tokenx2.TokenClaims{
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
	// hex-encoded HMAC-SHA256 is 64 chars
	if len(token) != 64 {
		t.Fatalf("expected 64-char hex token, got %d: %s", len(token), token)
	}
}

func TestHMACSigner_Sign_NonDeterministic(t *testing.T) {
	signer := tokenx2.NewHMACSigner("secret")
	claims := tokenx2.TokenClaims{
		UserID: "u", DeviceID: "d",
		IssuedAt:  time.Unix(100, 0),
		ExpiresAt: time.Unix(200, 0),
	}

	t1, _ := signer.Sign(context.Background(), claims)
	t2, _ := signer.Sign(context.Background(), claims)
	if t1 == t2 {
		t.Fatal("same claims should produce different tokens (random nonce)")
	}
}

func TestHMACSigner_Sign_DifferentClaimsProduceDifferentTokens(t *testing.T) {
	signer := tokenx2.NewHMACSigner("secret")
	claims1 := tokenx2.TokenClaims{
		UserID: "u", DeviceID: "d",
		IssuedAt: time.Unix(100, 0), ExpiresAt: time.Unix(200, 0),
	}
	claims2 := tokenx2.TokenClaims{
		UserID: "u", DeviceID: "d",
		IssuedAt: time.Unix(100, 0), ExpiresAt: time.Unix(300, 0),
	}

	t1, _ := signer.Sign(context.Background(), claims1)
	t2, _ := signer.Sign(context.Background(), claims2)
	if t1 == t2 {
		t.Fatal("different claims should produce different tokens")
	}
}

func TestHMACSigner_Verify_ReturnsNilNil(t *testing.T) {
	signer := tokenx2.NewHMACSigner("secret")
	claims, err := signer.Verify(context.Background(), "any-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claims != nil {
		t.Fatal("HMAC Verify should return nil claims (opaque token)")
	}
}
