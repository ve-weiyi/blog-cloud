package tokenx

import (
	"context"
	"fmt"
	"time"
)

const (
	TokenTypeBearer = "Bearer"
	TokenTypeSign   = "Sign"
)

type TokenHolder interface {
	TokenType() string
	VerifyToken(ctx context.Context, token string, uid string) error
	CreateToken(ctx context.Context, uid string, expires time.Duration) (string, error)
	RemoveToken(ctx context.Context, uid string) error
}

var (
	ErrTokenEmpty   = fmt.Errorf("token is empty")
	ErrTokenInvalid = fmt.Errorf("token is invalid")
	ErrTokenExpired = fmt.Errorf("token is expired")
)
