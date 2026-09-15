package tokenx_test

import (
	"context"
	"testing"
	"time"

	tokenx2 "github.com/ve-weiyi/blog-cloud/infra/tokenx"
)

func TestJWTSigner_Sign_ReturnsJWT(t *testing.T) {
	signer := tokenx2.NewJWTSigner([]byte("my-secret-key"), "test-issuer")
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
	signer := tokenx2.NewJWTSigner([]byte("secret"), "iss")
	claims := tokenx2.TokenClaims{
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
	signer := tokenx2.NewJWTSigner([]byte("secret"), "iss")
	claims := tokenx2.TokenClaims{
		UserID: "u", DeviceID: "d",
		IssuedAt:  time.Now().Add(-2 * time.Hour),
		ExpiresAt: time.Now().Add(-1 * time.Hour), // expired
	}

	token, err := signer.Sign(context.Background(), claims)
	if err != nil {
		t.Fatal(err)
	}

	_, err = signer.Verify(context.Background(), token)
	if err != tokenx2.ErrTokenExpired {
		t.Fatalf("expected ErrTokenExpired, got %v", err)
	}
}

func TestJWTSigner_Verify_InvalidToken(t *testing.T) {
	signer := tokenx2.NewJWTSigner([]byte("secret"), "iss")

	_, err := signer.Verify(context.Background(), "not.a.valid.jwt")
	if err != tokenx2.ErrTokenInvalid {
		t.Fatalf("expected ErrTokenInvalid, got %v", err)
	}
}

func TestJWTSigner_Verify_WrongKey(t *testing.T) {
	signer1 := tokenx2.NewJWTSigner([]byte("secret1"), "iss")
	signer2 := tokenx2.NewJWTSigner([]byte("secret2"), "iss")
	claims := tokenx2.TokenClaims{
		UserID: "u", DeviceID: "d",
		IssuedAt: time.Now(), ExpiresAt: time.Now().Add(time.Hour),
	}

	token, _ := signer1.Sign(context.Background(), claims)
	_, err := signer2.Verify(context.Background(), token)
	if err != tokenx2.ErrTokenInvalid {
		t.Fatalf("expected ErrTokenInvalid for wrong key, got %v", err)
	}
}
