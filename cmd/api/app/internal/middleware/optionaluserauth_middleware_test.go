package middleware

import (
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/ve-weiyi/blog-cloud/infra/http/headerx"
	"github.com/ve-weiyi/vkit/x/cryptox"
)

const testDeviceID = "0123456789abcdef0123456789abcdef"

func TestValidateGuestDevice(t *testing.T) {
	tests := []struct {
		name     string
		device   string
		stamp    int64
		token    func(device, stamp string) string
		wantFail bool
	}{
		{
			name:   "valid",
			device: testDeviceID,
			stamp:  time.Now().Unix(),
			token:  func(device, stamp string) string { return cryptox.Sha256v(device, stamp) },
		},
		{
			name:     "expired timestamp",
			device:   testDeviceID,
			stamp:    time.Now().Add(-time.Hour).Unix(),
			token:    func(device, stamp string) string { return cryptox.Sha256v(device, stamp) },
			wantFail: true,
		},
		{
			name:     "bad signature",
			device:   testDeviceID,
			stamp:    time.Now().Unix(),
			token:    func(string, string) string { return "bad" },
			wantFail: true,
		},
		{
			name:     "missing device",
			device:   "",
			stamp:    time.Now().Unix(),
			token:    func(string, string) string { return "bad" },
			wantFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/v1/ai/chat/stream", nil)
			stamp := strconv.FormatInt(tt.stamp, 10)
			if tt.device != "" {
				req.Header.Set(headerx.HeaderXDeviceId, tt.device)
			}
			req.Header.Set(headerx.HeaderTimestamp, stamp)
			req.Header.Set(headerx.HeaderXDeviceToken, tt.token(tt.device, stamp))

			err := validateGuestDevice(req, tt.device)
			if tt.wantFail && err == nil {
				t.Fatal("expected validation failure")
			}
			if !tt.wantFail && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
