package tokenx

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type hmacSigner struct {
	secret []byte
}

// NewHMACSigner 创建 HMAC-SHA256 签名器，生成不透明令牌。
func NewHMACSigner(secret string) Signer {
	return &hmacSigner{secret: []byte(secret)}
}

func (s *hmacSigner) Sign(_ context.Context, claims TokenClaims) (string, error) {
	// token = hex(HMAC-SHA256(key=secret, input=userId:deviceId:tokenType:issuedAtUnix:expiresAtUnix:random))
	nonce := make([]byte, 8)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	input := fmt.Sprintf("%s:%s:%s:%d:%d:%s",
		claims.UserID,
		claims.DeviceID,
		claims.TokenType,
		claims.IssuedAt.Unix(),
		claims.ExpiresAt.Unix(),
		hex.EncodeToString(nonce),
	)
	h := hmac.New(sha256.New, s.secret)
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Verify HMAC 令牌不可逆，始终返回 (nil, nil)。
func (s *hmacSigner) Verify(_ context.Context, _ string) (*TokenClaims, error) {
	return nil, nil
}
