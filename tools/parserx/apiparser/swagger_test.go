package apiparser

import (
	"fmt"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/swaggo/swag"

	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/jsonconv"
)

func Test_Load(t *testing.T) {
	// Example with default loaders defined at the package level
	doc, err := loads.Spec(SWAGER_PATH)
	if err != nil {
		fmt.Println("Could not load this spec")
		return
	}

	sp := doc.Spec()

	t.Log(jsonconv.AnyToJsonIndent(sp))
}

func Test_ParseAst(t *testing.T) {
	p := swag.New()

	p.ParseAPIMultiSearchDir([]string{"../../../blog-gin"}, "../../../blog-gin/main.go", 2)

	swagger := p.GetSwagger()
	fmt.Printf("Spec loaded: %v\n", jsonconv.AnyToJsonIndent(swagger))
}
