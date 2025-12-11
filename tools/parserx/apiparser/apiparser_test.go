package apiparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
)

const API_PATH = "../../../blog-gozero/service/api/blog/proto/blog.api"
const SWAGER_PATH = "../../../blog-gozero/service/api/blog/docs/blog.json"

func Test_NewApiParser(t *testing.T) {
	sp, err := parser.Parse(API_PATH)
	assert.Equal(t, nil, err)

	t.Log(jsonconv.AnyToJsonIndent(sp))
}

func Test_NewSpecParser(t *testing.T) {
	spa, err := NewSpecParser().ParseApi(API_PATH)
	assert.Equal(t, nil, err)

	t.Log(jsonconv.AnyToJsonIndent(spa.Service))
}

func Test_NewSwaggerParser(t *testing.T) {
	api, err := NewSwaggerParser().ParseApi(SWAGER_PATH)
	assert.Equal(t, nil, err)

	t.Log(jsonconv.AnyToJsonIndent(api.Service))
}

func Test_EqualParser(t *testing.T) {
	spa, err := NewSpecParser().ParseApi(API_PATH)
	assert.Equal(t, nil, err)
	//t.Log(jsonconv.AnyToJsonIndent(spa.Service))

	swa, err := NewSwaggerParser().ParseApi(SWAGER_PATH)
	assert.Equal(t, nil, err)
	//t.Log(jsonconv.AnyToJsonIndent(swa.Service))

	assert.Equal(t, spa.Types, swa.Types)
}
