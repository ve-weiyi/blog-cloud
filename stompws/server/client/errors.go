package client

import "fmt"

// Error messages following STOMP protocol conventions
type stompError string

func (e stompError) Error() string {
	return string(e)
}

// Standard STOMP error messages
const (
	errNotConnected         = stompError("expected CONNECT or STOMP frame")
	errAlreadyConnected     = stompError("already connected")
	errUnknownCommand       = stompError("unknown command")
	errSubscriptionExists   = stompError("subscription already exists")
	errSubscriptionNotFound = stompError("subscription not found")
	errTransactionExists    = stompError("transaction already exists")
	errTransactionNotFound  = stompError("unknown transaction")
	errInvalidHeartBeat     = stompError("invalid heart-beat format")
	errUnsupportedVersion   = stompError("unsupported protocol version")
	errAuthenticationFailed = stompError("authentication failed")
	errInvalidDestination   = stompError("invalid destination")
	errPermissionDenied     = stompError("permission denied")
)

func errMissingHeader(name string) error {
	return fmt.Errorf("missing required header: %s", name)
}
