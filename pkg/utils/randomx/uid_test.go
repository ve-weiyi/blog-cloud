package random

import (
	"testing"
)

func TestGenerateQQNumber(t *testing.T) {
	for i := 0; i < 10; i++ {
		qq := GenerateQQNumber()
		t.Logf("Generated QQ: %s", qq)
		if len(qq) < 6 || len(qq) > 8 {
			t.Errorf("QQ length is not in range [6, 8]: %s", qq)
		}
		if qq[0] == '0' {
			t.Errorf("QQ should not start with 0: %s", qq)
		}
	}
}
