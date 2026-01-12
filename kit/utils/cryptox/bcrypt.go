package cryptox

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// BcryptHash 使用 bcrypt 对密码进行加密,密码必须小于72位
func BcryptHash(password string) string {
	//bcrypt.DefaultCost默认数值10，编码一次100ms以内，可增大数值，增加破解时间成本，例如设置为14，编码一次1s以上
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	return string(bytes)
}

// BcryptCheck 对比明文密码和数据库的哈希值
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
