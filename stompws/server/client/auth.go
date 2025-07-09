package client

import (
	"fmt"
	"time"

	"github.com/go-stomp/stomp/v3/frame"
)

// authenticationFailed error with delay to prevent brute force
const authFailedDelay = time.Second

// Authenticator interface for authentication
type Authenticator interface {
	Authenticate(c *Client, f *frame.Frame) (clientId string, login string, err error)
}

// NoAuthenticator allows all connections
type NoAuthenticator struct{}

func NewNoAuthenticator() *NoAuthenticator {
	return &NoAuthenticator{}
}

func (a *NoAuthenticator) Authenticate(c *Client, f *frame.Frame) (string, string, error) {
	login := f.Header.Get(frame.Login)
	clientId := c.conn.RemoteAddr().String()

	if login == "" {
		login = clientId
	}
	return clientId, login, nil
}

// PasswordAuthenticator implements basic authentication
type PasswordAuthenticator struct {
	users map[string]string
}

// NewPasswordAuthenticator creates a new simple authenticator
func NewPasswordAuthenticator() *PasswordAuthenticator {
	return &PasswordAuthenticator{
		users: map[string]string{
			"guest": "guest",
			"admin": "admin",
		},
	}
}

func (a *PasswordAuthenticator) Authenticate(c *Client, f *frame.Frame) (string, string, error) {
	login := f.Header.Get(frame.Login)
	passcode := f.Header.Get(frame.Passcode)
	clientId := c.conn.RemoteAddr().String()

	if login == "" {
		return clientId, clientId, nil // Allow anonymous
	}

	expectedPasscode, exists := a.users[login]
	if !exists || expectedPasscode != passcode {
		return "", "", fmt.Errorf("authentication failed")
	}

	return clientId, login, nil
}

// AddUser adds a user to the authenticator
func (a *PasswordAuthenticator) AddUser(login, passcode string) {
	a.users[login] = passcode
}
