package stomphook

import (
	"fmt"

	"github.com/go-stomp/stomp/v3/frame"

	"github.com/ve-weiyi/stompws/server/client"
	"github.com/ve-weiyi/vkit/adapter/storex/tokenstore"
)

type SignAuthenticator struct {
	verifier tokenstore.TokenStore
}

func NewSignAuthenticator(verifier tokenstore.TokenStore) *SignAuthenticator {
	return &SignAuthenticator{
		verifier: verifier,
	}
}

// Authenticate implements the Authenticator interface
func (a *SignAuthenticator) Authenticate(c *client.Client, f *frame.Frame) (string, string, error) {
	login := f.Header.Get("login")
	passcode := f.Header.Get("passcode")
	clientId := f.Header.Get("client")

	if clientId == "" {
		return "", "", fmt.Errorf("stomp auth failed: missing header: 'client'")
	}

	// 游客模式
	if login == "" && passcode == "" {
		return clientId, login, nil
	}

	// token校验
	err := a.verifier.ValidateToken(login, passcode)
	if err != nil {
		return "", "", fmt.Errorf("stomp auth failed: %v", err)
	}

	return clientId, login, nil
}
