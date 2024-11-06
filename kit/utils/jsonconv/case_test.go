package jsonconv

import (
	"log"
	"testing"
)

func TestCamel2Case(t *testing.T) {
	str := "link__intro"
	log.Println("--->", str)

	// 转下划线
	cases := Case2Snake(str)
	log.Println("--->", cases)

	// 转驼峰
	camel := Case2Camel(cases)
	log.Println("--->", camel)

	// 转下划线
	cases = Case2Snake(camel)
	log.Println("--->", cases)
}
