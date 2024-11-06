package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5v md5 加盐加密
func Md5v(str string, salt string) string {
	h := md5.New()
	h.Write([]byte(salt + str))
	var res []byte
	res = h.Sum(nil)
	return hex.EncodeToString(res)
}

// Md5vMulti  iteration:加密次数
func Md5vMulti(str string, salt string, iteration int) string {
	h := md5.New()
	h.Write([]byte(salt + str))
	var res []byte
	res = h.Sum(nil)
	for i := 0; i < iteration-1; i++ {
		h.Reset()
		h.Write(res)
		res = h.Sum(nil)
	}
	return hex.EncodeToString(res)
}
