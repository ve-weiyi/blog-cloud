package parserx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
)

func Test_NewApiParser(t *testing.T) {
	sp, err := parser.Parse("/Users/weiyi/Github/ve-blog-golang/zero/service/api/blog/proto/blog.api")
	assert.Equal(t, nil, err)

	t.Log(jsonconv.AnyToJsonIndent(sp))
}

func Test_NewSpecParser(t *testing.T) {

	spa, err := NewSpecParser().ParseApi("../../testdata/test.api")
	assert.Equal(t, nil, err)

	t.Log(jsonconv.AnyToJsonIndent(spa.Service))

	//swa, err := NewSwaggerParser().ParseApi("../testdata/test.json")
	//assert.Equal(t, nil, err)
	//
	//t.Log(jsonconv.AnyToJsonIndent(swa.Service))
}
