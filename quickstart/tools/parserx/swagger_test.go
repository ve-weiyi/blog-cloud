package parserx

import (
	"fmt"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/stretchr/testify/assert"
	"github.com/swaggo/swag"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
)

func Test_NewSwaggerParser(t *testing.T) {
	api, err := NewSwaggerParser().ParseApi("./_testdata/test.json")
	assert.Equal(t, nil, err)

	t.Log(jsonconv.AnyToJsonIndent(api.Types))
}

func Test_Load(t *testing.T) {
	// Example with default loaders defined at the package level

	path := "./_testdata/test.json"
	doc, err := loads.Spec(path)
	if err != nil {
		fmt.Println("Could not load this spec")
		return
	}

	sp := doc.Spec()

	t.Log(jsonconv.AnyToJsonIndent(sp))
}

func Test_ParseAst(t *testing.T) {
	p := swag.New()

	p.ParseAPIMultiSearchDir([]string{"/Users/weiyi/Github/ve-blog-golang/server/api/blog/controller"}, "/Users/weiyi/Github/ve-blog-golang/server/main.go", 2)

	swagger := p.GetSwagger()
	fmt.Printf("Spec loaded: %v\n", jsonconv.AnyToJsonIndent(swagger))
}
