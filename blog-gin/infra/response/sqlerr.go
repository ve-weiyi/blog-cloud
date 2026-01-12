package response

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	"github.com/go-sql-driver/mysql"
)

// mysql 错误码
const (
	SqlKeyDuplicateCode = 1062
	DataTooLongError    = 1406
)

// https://fromdual.com/mysql-error-codes-and-messages-1050-1099
func SqlErrorI18n(errorMessage *mysql.MySQLError) string {
	var message string
	switch int(errorMessage.Number) {
	case SqlKeyDuplicateCode:
		// Message: Duplicate entry '%s' for key %d
		reg := regexp.MustCompile(`Duplicate entry ('(.*?)')`)
		subText := reg.FindString(errorMessage.Message)
		start := len("Duplicate entry '")
		msg := []rune(subText)[start : utf8.RuneCountInString(subText)-1]
		message = fmt.Sprintf("%s已存在，请更改后提交", string(msg))
	case DataTooLongError:
		// Data too long for column 'value' at row 1
		reg := regexp.MustCompile(`Data too long for column ('(.*?)')`)
		subText := reg.FindString(errorMessage.Error())
		start := len("Data too long for column '")
		column := []rune(subText)[start : utf8.RuneCountInString(subText)-1]
		msg := "tableColumns." + string(column)
		message = fmt.Sprintf("%s 内容超过数据库字段限制，请重新输入", msg)
	default:
		message = errorMessage.Error()
	}
	return message
}
