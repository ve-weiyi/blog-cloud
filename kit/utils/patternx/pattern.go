package patternx

import "regexp"

// 使用正则表达式验证邮箱地址的格式
func IsValidEmail(email string) bool {
	// 正则表达式模式
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// 编译正则表达式
	reg := regexp.MustCompile(pattern)

	// 匹配邮箱地址
	return reg.MatchString(email)
}

// 使用正则表达式验证手机号的格式
func IsValidPhone(photo string) bool {
	pattern := `^1[3456789]\d{9}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(photo)
}

// 全数字
func IsValidDigit(digit string) bool {
	pattern := `^\d+$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(digit)
}
