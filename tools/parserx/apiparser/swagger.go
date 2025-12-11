package apiparser

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"

	"github.com/ve-weiyi/ve-blog-golang/tools/parserx/apiparser/aspec"
)

type SwaggerOption func(*SwaggerParser)

type SwaggerParser struct {
}

func NewSwaggerParser(opts ...SwaggerOption) ApiParser {
	return &SwaggerParser{}
}

func (s *SwaggerParser) ParseApi(filename string) (out *aspec.ApiSpec, err error) {
	swagger, err := ParseSwagJson(filename)
	if err != nil {
		return nil, err
	}

	gps := getGroupsFromSwag(swagger)

	tps := getTypesFromSwag(swagger)

	out = &aspec.ApiSpec{
		Service: aspec.Service{
			Name:   swagger.Host,
			Groups: gps,
		},
		Types: tps,
	}

	return out, nil
}

func ParseSwagJson(filename string) (out *spec.Swagger, err error) {
	doc, err := loads.Spec(filename)
	if err != nil {
		return nil, err
	}

	sp := doc.Spec()
	return sp, nil
}

func getGroupsFromSwag(sp *spec.Swagger) []aspec.Group {
	var gps []aspec.Group
	for k, v := range sp.Paths.Paths {
		var rts []aspec.Route

		var mr map[string]spec.Operation
		jb, _ := v.MarshalJSON()
		_ = json.Unmarshal(jb, &mr)

		for m, r := range mr {
			var h, q, p, f []aspec.Member
			var req, resp aspec.DefineStruct
			for _, parameter := range r.Parameters {
				switch parameter.In {
				case "header":
					h = append(h, aspec.Member{
						Name: parameter.Name,
						Type: aspec.PrimitiveType{
							RawName: parameter.Type,
						},
					})
				case "query":
					q = append(q, aspec.Member{
						Name: parameter.Name,
						Type: aspec.PrimitiveType{
							RawName: parameter.Type,
						},
					})
				case "path":
					p = append(p, aspec.Member{
						Name: parameter.Name,
						Type: aspec.PrimitiveType{
							RawName: parameter.Type,
						},
					})
				case "form":
					f = append(f, aspec.Member{
						Name: parameter.Name,
						Type: aspec.PrimitiveType{
							RawName: parameter.Type,
						},
					})
				case "body":
					if parameter.Schema != nil {
						req.RawName = strings.TrimPrefix(parameter.Schema.Ref.String(), "#/definitions/")
					}
				default:
				}
			}

			if r.Responses != nil {
				for _, response := range r.Responses.StatusCodeResponses {
					if response.Schema != nil {
						data := response.Schema.Properties["data"]
						resp.RawName = strings.TrimPrefix(data.Ref.String(), "#/definitions/")
					}
				}
			}

			rt := aspec.Route{
				Method:         m,
				Path:           k,
				PathType:       p,
				HeaderType:     h,
				FormType:       f,
				QueryType:      q,
				RequestType:    req,
				ResponseType:   resp,
				Docs:           nil,
				Handler:        r.ID,
				AtDoc:          aspec.AtDoc{},
				HandlerDoc:     nil,
				HandlerComment: nil,
				Doc:            []string{r.Summary},
				Comment:        nil,
			}

			rts = append(rts, rt)
		}

		gp := aspec.Group{
			Annotation: aspec.Annotation{
				Properties: map[string]string{
					"group": "",
				},
			},
			Routes: rts,
		}

		gps = append(gps, gp)
	}
	return gps
}

func getTypesFromSwag(sp *spec.Swagger) []aspec.Type {
	var types []aspec.Type

	for k, v := range sp.Definitions {
		switch v.Type {
		default:
			var ms []aspec.Member
			for n, m := range v.Properties {
				ms = append(ms, aspec.Member{
					Name:     n,
					Type:     convertSchema(m),
					Tag:      fmt.Sprintf("`json:%v`", n),
					Comment:  fmt.Sprintf("// %v", m.Description),
					Docs:     make(aspec.Doc, 0),
					IsInline: false,
				})
			}

			model := aspec.DefineStruct{
				RawName: k,
				Members: ms,
				Docs:    nil,
			}

			types = append(types, model)
		}

	}

	sort.Slice(types, func(i, j int) bool {
		return types[i].Name() < types[j].Name()
	})

	return types
}

func convertSchema(in spec.Schema) aspec.Type {
	if ref := in.Ref.String(); ref != "" {
		if len(ref) > 14 && ref[:14] == "#/definitions/" {
			ref = ref[14:]
		}
		return aspec.DefineStruct{RawName: ref}
	}

	if len(in.Type) != 1 {
		return aspec.InterfaceType{RawName: "interface{}"}
	}

	switch in.Type[0] {
	case "array":
		if in.Items != nil && in.Items.Schema != nil {
			elemType := convertSchema(*in.Items.Schema)
			return aspec.ArrayType{
				RawName: "[]" + elemType.Name(),
				Value:   elemType,
			}
		}
		return aspec.ArrayType{RawName: "[]interface{}"}

	case "object":
		if len(in.Properties) > 0 {
			var ms []aspec.Member
			for n, m := range in.Properties {
				ms = append(ms, aspec.Member{
					Name: n,
					Type: convertSchema(m),
				})
			}
			return aspec.DefineStruct{Members: ms}
		}
		return aspec.InterfaceType{RawName: "interface{}"}

	default:
		return aspec.PrimitiveType{RawName: convertSwagTypeToGoType(in.Type[0])}
	}
}

func convertSwagTypeToGoType(t string) string {
	switch t {
	case "integer":
		return "int64"
	default:
		return t
	}
}
