package mark

import (
	"regexp"
)

// 手机号脱敏 13812345678-->138****5678
func maskPhone(phone string) string {
	re := regexp.MustCompile(`(\d{3})\d{4}(\d{4})`)
	masked := re.ReplaceAllString(phone, "$1****$2")
	return masked
}

// 邮箱脱敏 username@qq.com-->us****@qq.com
func maskEmail(email string) string {
	re := regexp.MustCompile(`([\w._%+-]+)(@[\w.-]+\.[A-Za-z]{2,})`)
	masked := re.ReplaceAllStringFunc(email, func(s string) string {
		matches := re.FindStringSubmatch(s)
		if len(matches) == 3 {
			username := matches[1]
			domain := matches[2]

			//替换username-->us****
			runes := []rune(username)
			length := len(runes)
			if length <= 2 {
				return string(runes)
			}
			masked := make([]rune, length)
			masked[0] = runes[0]
			masked[1] = runes[2]
			for i := 2; i < length; i++ {
				masked[i] = '*'
			}
			//返回us***@email.com
			return string(masked) + domain
		}
		return s
	})
	return masked
}
