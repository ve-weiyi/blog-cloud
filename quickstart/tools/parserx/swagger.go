package parserx

import (
	"encoding/json"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"

	"github.com/ve-weiyi/ve-blog-golang/quickstart/tools/parserx/aspec"
)

type SwaggerOption func(*SwaggerParser)

type SwaggerParser struct {
	Base string
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
			Name:   swagger.Info.Title,
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
						req.RawName = parameter.Schema.Ref.String()
					}
				default:
				}
			}

			if r.Responses != nil {
				for _, response := range r.Responses.StatusCodeResponses {
					if response.Schema != nil {
						resp.RawName = response.Schema.Ref.String()
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
					Tag:      "",
					Comment:  "",
					Docs:     make(aspec.Doc, 0),
					IsInline: false,
				})
			}

			model := &aspec.DefineStruct{
				RawName: k,
				Members: ms,
				Docs:    nil,
			}

			types = append(types, model)
		}

	}

	return types
}

func convertSchema(in spec.Schema) aspec.Type {
	var tp aspec.Type
	if len(in.Type) == 1 {
		it := in.Type[0]
		switch it {
		case "object":

		default: // string、int
			tp = aspec.PrimitiveType{
				RawName: it,
			}
		}

	}

	return tp
}
