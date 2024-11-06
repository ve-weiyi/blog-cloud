package invent

import (
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
)

var StdMapUtils = map[string]any{
	"Case2Camel": jsonconv.Case2Camel,
	"Case2Snake": jsonconv.Case2Snake,
	"ToUpper":    strings.ToUpper,
	"ToLower":    strings.ToLower,
}
