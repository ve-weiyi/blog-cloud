package stomphook

import (
	"fmt"

	"github.com/go-stomp/stomp/v3/frame"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/infra/tokenx"
	"github.com/ve-weiyi/ve-blog-golang/stompws/server/client"
)

type JwtAuthenticator struct {
	verifier tokenx.TokenManager
}

func NewJwtAuthenticator(verifier tokenx.TokenManager) *JwtAuthenticator {
	return &JwtAuthenticator{
		verifier: verifier,
	}
}

// Authenticate implements the Authenticator interface
func (a *JwtAuthenticator) Authenticate(c *client.Client, f *frame.Frame) (string, string, error) {
	login := f.Header.Get("login")
	passcode := f.Header.Get("passcode")
	clientId := f.Header.Get("client")

	if login == "" {
		return "", "", fmt.Errorf("stomp auth failed: missing header: 'login'")
	}
	if passcode == "" {
		return "", "", fmt.Errorf("stomp auth failed: missing header: 'passcode'")
	}
	if clientId == "" {
		return "", "", fmt.Errorf("stomp auth failed: missing header: 'client'")
	}
	// 校验jwt
	err := a.verifier.ValidateToken(login, passcode)
	if err != nil {
		return "", "", fmt.Errorf("stomp auth failed: %v", err)
	}

	return clientId, login, nil
}
