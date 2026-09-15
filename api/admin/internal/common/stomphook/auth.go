package stomphook

import (
	"context"
	"fmt"

	"github.com/go-stomp/stomp/v3/frame"

	"github.com/ve-weiyi/stompws/server/client"

	"github.com/ve-weiyi/blog-cloud/infra/tokenx"
)

type JwtAuthenticator struct {
	store tokenx.Manager
}

func NewJwtAuthenticator(store tokenx.Manager) *JwtAuthenticator {
	return &JwtAuthenticator{
		store: store,
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
	err := a.store.Validate(context.Background(), login, clientId, passcode)
	if err != nil {
		return "", "", fmt.Errorf("stomp auth failed: %v", err)
	}

	return clientId, login, nil
}
