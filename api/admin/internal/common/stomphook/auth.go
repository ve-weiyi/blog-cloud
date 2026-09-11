package stomphook

import (
	"fmt"

	"github.com/go-stomp/stomp/v3/frame"

	"github.com/ve-weiyi/stompws/server/client"
	"github.com/ve-weiyi/vkit/adapter/storex/tokenstore"
)

type JwtAuthenticator struct {
	store tokenstore.TokenStore
}

func NewJwtAuthenticator(store tokenstore.TokenStore) *JwtAuthenticator {
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
	err := a.store.ValidateToken(login, passcode)
	if err != nil {
		return "", "", fmt.Errorf("stomp auth failed: %v", err)
	}

	return clientId, login, nil
}
