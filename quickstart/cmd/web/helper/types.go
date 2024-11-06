package helper

import (
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/convertx"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
	"github.com/ve-weiyi/ve-blog-golang/quickstart/tools/parserx/aspec"
)

type TsApiGroup struct {
	Name       string // account
	Prefix     string
	Middleware []string
	Routes     []TsApiRoute
}

type TsApiRoute struct {
	Summery  string
	Handler  string
	Path     string
	Method   string
	Request  string
	Response string
}

type TsType struct {
	Comment string
	Name    string
	Extends []string
	Fields  []TsTypeField
}

type TsTypeField struct {
	Name     string
	Type     string
	Comment  string
	Nullable bool
}

func ConvertApiTs(sp *aspec.ApiSpec) map[string][]TsApiGroup {
	var groups []TsApiGroup
	for _, v := range sp.Service.Groups {
		var routes []TsApiRoute
		for _, r := range v.Routes {

			var doc string
			for _, d := range r.Doc {
				doc = doc + d
			}

			for _, d := range r.HandlerDoc {
				doc = doc + d
			}

			if r.AtDoc.Text != "" {
				doc = doc + strings.Trim(strings.Trim(r.AtDoc.Text, "\\"), "\"")
			}

			req := "any"
			resp := "any"
			if r.RequestType != nil {
				req = convertx.ConvertGoTypeToTsType(r.RequestType.Name())
			}

			if r.ResponseType != nil {
				resp = convertx.ConvertGoTypeToTsType(r.ResponseType.Name())
			}

			rt := TsApiRoute{
				Summery:  doc,
				Handler:  jsonconv.FirstLower(r.Handler) + "Api",
				Path:     r.Path,
				Method:   strings.ToUpper(r.Method),
				Request:  req,
				Response: resp,
			}

			routes = append(routes, rt)
		}

		var name = v.Annotation.Properties["group"]
		if name == "" {
			name = "common"
		}

		var prefix = v.Annotation.Properties["prefix"]

		var middleware = v.Annotation.Properties["middleware"]

		g := TsApiGroup{
			Name:       name,
			Prefix:     prefix,
			Middleware: []string{},
			Routes:     routes,
		}

		if middleware != "" {
			g.Middleware = strings.Split(middleware, ",")
		}

		groups = append(groups, g)
	}

	mgs := make(map[string][]TsApiGroup)
	for _, v := range groups {
		mgs[v.Name] = append(mgs[v.Name], v)
	}

	return mgs
}

func ConvertTypeTs(sp *aspec.ApiSpec) map[string]TsType {
	ts := make(map[string]TsType)

	for _, v := range sp.Types {
		ts[v.Name()] = convertTypeTs(v)
	}

	return ts
}

func convertTypeTs(st aspec.Type) TsType {
	var ts TsType

	switch t := st.(type) {
	case aspec.DefineStruct:

		var ex []string
		var tfs []TsTypeField
		for _, v := range t.Members {
			if v.Name == "" {
				ex = append(ex, v.Type.Name())
			} else {
				m := TsTypeField{
					Comment:  v.Comment,
					Name:     jsonconv.Case2Snake(v.Name),
					Type:     convertx.ConvertGoTypeToTsType(v.Type.Name()),
					Nullable: strings.HasPrefix(v.Type.Name(), "*") || strings.Contains(v.Tag, "optional"),
				}
				tfs = append(tfs, m)
			}
		}

		ts = TsType{
			Comment: strings.Join(t.Comments(), "\n"),
			Name:    t.Name(),
			Extends: ex,
			Fields:  tfs,
		}

	case aspec.PrimitiveType:
	case aspec.MapType:
	case aspec.ArrayType:
	case aspec.InterfaceType:
	case aspec.PointerType:
	default:
		panic("unknown type")
	}

	return ts
}
