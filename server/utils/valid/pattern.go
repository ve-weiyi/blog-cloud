package valid

import "regexp"

// 使用正则表达式验证邮箱地址的格式
func IsEmailValid(email string) bool {
	// 正则表达式模式
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// 编译正则表达式
	reg := regexp.MustCompile(pattern)

	// 匹配邮箱地址
	return reg.MatchString(email)
}
