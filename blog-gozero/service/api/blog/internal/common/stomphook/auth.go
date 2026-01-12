package stomphook

import (
	"context"
	"fmt"

	"github.com/go-stomp/stomp/v3/frame"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/infra/tokenx"
	"github.com/ve-weiyi/ve-blog-golang/stompws/server/client"
)

type SignAuthenticator struct {
	verifier tokenx.TokenHolder
}

func NewSignAuthenticator(verifier tokenx.TokenHolder) *SignAuthenticator {
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
	err := a.verifier.VerifyToken(context.Background(), passcode, login)
	if err != nil {
		return "", "", fmt.Errorf("stomp auth failed: %v", err)
	}

	return clientId, login, nil
}
