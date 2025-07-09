package docs

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
)

//go:embed admin.json
var docs string

func Test_Load(t *testing.T) {

	doc, err := loads.Analyzed(json.RawMessage(docs), "")
	if err != nil {
		fmt.Println("Could not load this spec")
		return
	}

	sp := doc.Spec()

	//for k, v := range sp.Paths.Paths {
	//	t.Logf("path: %s,method: %s", k, jsonconv.AnyToJsonIndent(v))
	//}
	t.Log(jsonconv.AnyToJsonIndent(sp))

	routes := getRoutes(sp)

	for _, v := range routes {
		for _, o := range v {
			if o != nil {

			}
		}
	}

	//var out map[string]map[string]any
	//jsonconv.AnyToAny(sp.Paths.Paths, &out)
	//t.Log(jsonconv.AnyToJsonIndent(out))
}

func getRoutes(sp *spec.Swagger) map[string]map[string]*spec.Operation {
	// map[path][method] -> operation
	routes := make(map[string]map[string]*spec.Operation)

	for k, v := range sp.Paths.Paths {
		if routes[k] == nil {
			routes[k] = make(map[string]*spec.Operation)
		}

		if v.Get != nil {
			routes[k][http.MethodGet] = v.Get
		}

		if v.Put != nil {
			routes[k][http.MethodPut] = v.Put
		}

		if v.Post != nil {
			routes[k][http.MethodPost] = v.Post
		}

		if v.Delete != nil {
			routes[k][http.MethodDelete] = v.Delete
		}

		if v.Options != nil {
			routes[k][http.MethodOptions] = v.Options
		}

		if v.Head != nil {
			routes[k][http.MethodHead] = v.Head
		}
	}

	return routes
}
