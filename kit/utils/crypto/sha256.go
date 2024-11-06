package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

// Md5v md5 加盐加密
func Sha256v(str string, salt string) string {
	h := sha256.New()
	h.Write([]byte(salt + str))
	var res []byte
	res = h.Sum(nil)
	return hex.EncodeToString(res)
}
