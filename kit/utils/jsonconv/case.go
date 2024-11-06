package jsonconv

import (
	"strings"
	"unicode"
)

const (
	Camel = 0 //驼峰
	Case  = 1 //下划线
)

var replaceKey = map[string]string{
	//"DTO": "dto",
	//"UID": "uid",
	//"PID": "pid",
	//"ID":  "id",
	//"URL": "url",
	//"AI":  "ai",
}

/**
 * 驼峰式写法转为下划线写法
 * @description XxYx->xx_yy	XxYY->xx_yy	EncodeURL->url  TagDTOList->tag_dto_list
 **/
func Case2Snake(XxYY string) string {
	xx_y_y := make([]byte, 0)
	i := 0

	for i < len(XxYY) {
		/** 替换不需要转换的字符大小写 **/
		found := false
		for prefix, replace := range replaceKey {
			if strings.HasPrefix(XxYY[i:], prefix) {
				// 非首个字符
				if len(xx_y_y) != 0 {
					xx_y_y = append(xx_y_y, '_')
				}

				xx_y_y = append(xx_y_y, []byte(replace)...)
				i += len(prefix)
				found = true
				break
			}
		}
		if found {
			continue
		}
		/**  **/

		flag := false
		// 非首个字符
		if len(xx_y_y) != 0 {
			// 方案二:自动匹配。如果前一个字符是小写 || 后一个字符是小写
			//if (0 < i-1 && unicode.IsLower(rune(XxYY[i-1]))) || (i+1 < len(XxYY) && unicode.IsLower(rune(XxYY[i+1]))) {
			//	xx_y_y = append(xx_y_y, '_')
			//}
			flag = true
		}

		// 未找到 replaceKey ，进行正常转换
		w := rune(XxYY[i])
		i++
		// 遇到数字
		if unicode.IsDigit(w) {
			xx_y_y = append(xx_y_y, byte(w))
			continue
		}
		// 遇到非字母
		if !unicode.IsLetter(w) {
			xx_y_y = append(xx_y_y, byte('_'))
			continue
		}
		// 如果是大写
		if unicode.IsUpper(w) {
			if flag {
				xx_y_y = append(xx_y_y, '_')
			}
			xx_y_y = append(xx_y_y, byte(unicode.ToLower(w)))
		} else {
			xx_y_y = append(xx_y_y, byte(w))
		}
	}

	return string(xx_y_y)
}

/**
 * 下划线转驼峰
 * @description xx_yy to XxYx  xx_y_y to XxYY  XxYY to XxYY
 * @date 2023/2/15
 * @param xx_y_y
 * @return XxYY
 **/
func Case2Camel(xx_y_y string) string {
	//id类型转换大写
	XxYY := make([]byte, 0, len(xx_y_y))
	//是否遇到下划线,初始化值为true则转换第一个字母
	line := true
	i := 0
	for i < len(xx_y_y) {
		/** 替换不需要转换的字符大小写 **/
		found := false
		for key, value := range replaceKey {
			if strings.HasPrefix(xx_y_y[i:], value) && len(key) == len(xx_y_y[i:]) {
				// 非首个字符
				XxYY = append(XxYY, []byte(key)...)
				i += len(value)
				found = true
				break
			}
		}
		if found {
			continue
		}
		/**  **/

		// 未找到 replaceKey ，进行正常转换
		w := rune(xx_y_y[i])
		i++
		//遇到数字
		if unicode.IsDigit(w) {
			XxYY = append(XxYY, byte(w))
			continue
		}

		//遇到 _
		if !unicode.IsLetter(w) {
			line = true
			continue
		}

		//遇到小写
		if w >= 'a' && w <= 'z' {
			if line {
				w = w - 32
			}
		}
		//遇到大写，跳过
		if w >= 'A' && w <= 'Z' {

		}
		//只对 _ 后一个字母生效
		if line {
			line = false
		}
		XxYY = append(XxYY, byte(w))
	}

	return string(XxYY[:])
}
