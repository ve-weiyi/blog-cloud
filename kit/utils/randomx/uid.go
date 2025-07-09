package random

import (
	"math/rand"
	"time"
)

// 初始化随机种子
func init() {
	rand.Seed(time.Now().UnixNano())
}

// 生成随机数字账号，不以 0 开头
func GenerateQQNumber() string {
	length := rand.Intn(2) + 7 // 7~8 位
	digits := make([]byte, length)

	// 首位不能为 0
	digits[0] = byte(rand.Intn(9)+1) + '0'

	for i := 1; i < length; i++ {
		digits[i] = byte(rand.Intn(10)) + '0'
	}
	return string(digits)
}
