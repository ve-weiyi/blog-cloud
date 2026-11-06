package chat

import (
	"strings"
	"testing"
)

func TestNormalizeContent(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "trim", input: "  hello  ", want: "hello"},
		{name: "empty", input: "   ", wantErr: true},
		{name: "max", input: strings.Repeat("好", MaxContentRunes), want: strings.Repeat("好", MaxContentRunes)},
		{name: "too long", input: strings.Repeat("好", MaxContentRunes+1), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeContent(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValidType(t *testing.T) {
	if !ValidType(TypeText) {
		t.Fatal("text should be supported")
	}
	for _, messageType := range []string{"", "image", "file", "audio", "video"} {
		if ValidType(messageType) {
			t.Fatalf("%q should be rejected", messageType)
		}
	}
}
