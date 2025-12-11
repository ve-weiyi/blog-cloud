package stomphook

import (
	"context"
	"fmt"

	"github.com/go-stomp/stomp/v3/frame"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/tokenx"
	"github.com/ve-weiyi/ve-blog-golang/stompws/server/client"
)

type JwtAuthenticator struct {
	verifier tokenx.TokenHolder
}

func NewJwtAuthenticator(verifier tokenx.TokenHolder) *JwtAuthenticator {
	return &JwtAuthenticator{
		verifier: verifier,
	}
}

// Authenticate implements the Authenticator interface
func (a *JwtAuthenticator) Authenticate(c *client.Client, f *frame.Frame) (string, string, error) {
	login := f.Header.Get("login")
	token := f.Header.Get("authorization")

	if login == "" || token == "" {
		return "", "", fmt.Errorf("stomp auth failed: missing login or authorization header")
	}

	// 校验jwt
	err := a.verifier.VerifyToken(context.Background(), token, login)
	if err != nil {
		return "", "", fmt.Errorf("stomp auth failed: %v", err)
	}

	return login, login, nil
}
