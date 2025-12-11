package client

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-stomp/stomp/v3"
	"github.com/go-stomp/stomp/v3/frame"
)

var (
	heartBeatPattern = regexp.MustCompile("^[0-9]{1,9},[0-9]{1,9}$")
	destPattern      = regexp.MustCompile(`^/[a-zA-Z0-9/_-]+$`)
)

func validateHeartBeat(hb string) bool {
	if !heartBeatPattern.MatchString(hb) {
		return false
	}
	parts := strings.Split(hb, ",")
	cx, _ := strconv.ParseUint(parts[0], 10, 32)
	cy, _ := strconv.ParseUint(parts[1], 10, 32)
	return cx <= maxHeartBeat && cy <= maxHeartBeat
}

func validateDestination(dest string) bool {
	return destPattern.MatchString(dest)
}

type validator struct {
	version stomp.Version
}

func newValidator(version stomp.Version) stomp.Validator {
	return &validator{version: version}
}

func (v *validator) Validate(f *frame.Frame) error {
	switch f.Command {
	case frame.CONNECT, frame.STOMP:
		if _, ok := f.Header.Contains(frame.Receipt); ok {
			return fmt.Errorf("receipt header not allowed in CONNECT/STOMP frame")
		}
		if v.version != stomp.V10 {
			if _, ok := f.Header.Contains(frame.AcceptVersion); !ok {
				return errMissingHeader(frame.AcceptVersion)
			}
			if _, ok := f.Header.Contains(frame.Host); !ok {
				return errMissingHeader(frame.Host)
			}
		}
		if hb, ok := f.Header.Contains(frame.HeartBeat); ok && !validateHeartBeat(hb) {
			return errInvalidHeartBeat
		}
	case frame.SEND:
		if _, ok := f.Header.Contains(frame.Destination); !ok {
			return errMissingHeader(frame.Destination)
		}
		if !validateDestination(f.Header.Get(frame.Destination)) {
			return errInvalidDestination
		}
	case frame.SUBSCRIBE:
		if _, ok := f.Header.Contains(frame.Destination); !ok {
			return errMissingHeader(frame.Destination)
		}
		if !validateDestination(f.Header.Get(frame.Destination)) {
			return errInvalidDestination
		}
		if _, ok := f.Header.Contains(frame.Id); !ok {
			return errMissingHeader(frame.Id)
		}
	case frame.UNSUBSCRIBE:
		if _, ok := f.Header.Contains(frame.Id); !ok {
			return errMissingHeader(frame.Id)
		}
	case frame.ACK:
		if v.version == stomp.V12 {
			if _, ok := f.Header.Contains(frame.Id); !ok {
				return errMissingHeader(frame.Id)
			}
		} else if _, ok := f.Header.Contains(frame.MessageId); !ok {
			if _, ok := f.Header.Contains(frame.Ack); !ok {
				return errMissingHeader("message-id or ack")
			}
		}
	case frame.NACK:
		if v.version == stomp.V10 {
			return fmt.Errorf("NACK not supported in STOMP 1.0")
		}
		if v.version == stomp.V12 {
			if _, ok := f.Header.Contains(frame.Id); !ok {
				return errMissingHeader(frame.Id)
			}
		} else if _, ok := f.Header.Contains(frame.MessageId); !ok {
			if _, ok := f.Header.Contains(frame.Ack); !ok {
				return errMissingHeader("message-id or ack")
			}
		}
	case frame.BEGIN, frame.COMMIT, frame.ABORT:
		if _, ok := f.Header.Contains(frame.Transaction); !ok {
			return errMissingHeader(frame.Transaction)
		}
	}
	return nil
}
