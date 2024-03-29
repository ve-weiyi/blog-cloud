package mark

import (
	"fmt"
	"testing"
)

func TestMask(t *testing.T) {
	phone := "我的手机号是13812345678"
	email := "我的邮箱号是username@qq.com"

	maskedPhone := maskPhone(phone)
	maskedEmail := maskEmail(email)

	fmt.Println(maskedPhone)
	fmt.Println(maskedEmail)
}
