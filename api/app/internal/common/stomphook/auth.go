package stomphook

import (
	"context"
	"fmt"

	"github.com/go-stomp/stomp/v3/frame"

	"github.com/ve-weiyi/blog-cloud/infra/storex/tokenx"
	"github.com/ve-weiyi/stompws/server/client"
)

type SignAuthenticator struct {
	verifier tokenx.Manager
}

func NewSignAuthenticator(verifier tokenx.Manager) *SignAuthenticator {
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
	err := a.verifier.Validate(context.Background(), login, clientId, passcode)
	if err != nil {
		return "", "", fmt.Errorf("stomp auth failed: %v", err)
	}

	return clientId, login, nil
}
