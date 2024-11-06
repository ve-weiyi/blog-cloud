package helper

import (
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
	"github.com/ve-weiyi/ve-blog-golang/quickstart/tools/parserx/aspec"
)

type GroupRoute struct {
	Name       string // account
	Prefix     string
	Middleware []string
	Routes     []Route
}

type Route struct {
	Doc      string // 用户接口
	Handler  string // UserHandler
	Path     string // /api/v1/user
	Method   string // POST、GET、PUT、DELETE
	Request  string
	Response string
}

type GroupType struct {
	Group string
	Types []spec.Type
}

func ConvertRouteGroups(sp *aspec.ApiSpec) (out map[string][]GroupRoute) {
	var groups []GroupRoute
	for _, v := range sp.Service.Groups {
		var routes []Route
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

			var req, resp string
			if r.RequestType != nil {
				req = r.RequestType.Name()
			}

			if r.ResponseType != nil {
				resp = r.ResponseType.Name()
			}

			rt := Route{
				Method:   strings.ToUpper(r.Method),
				Path:     r.Path,
				Handler:  jsonconv.Case2Camel(r.Handler),
				Doc:      doc,
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

		g := GroupRoute{
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

	out = make(map[string][]GroupRoute)
	for _, v := range groups {
		out[v.Name] = append(out[v.Name], v)
	}

	return out
}

func ConvertTypeGroups(sp *spec.ApiSpec) (out []GroupType) {
	var groups []GroupType

	mt := make(map[string]spec.Type)
	for _, v := range sp.Types {
		mt[v.Name()] = v
	}

	tps := make(map[string][]spec.Type)

	for _, v := range sp.Service.Groups {
		for _, r := range v.Routes {
			group := v.Annotation.Properties["group"]
			if group == "" {
				group = "common"
			}
			if r.RequestType != nil {
				if mt[r.RequestType.Name()] != nil {
					tps[group] = append(tps[group], mt[r.RequestType.Name()])
				}
			}

			if r.ResponseType != nil {
				if mt[r.ResponseType.Name()] != nil {
					tps[group] = append(tps[group], mt[r.ResponseType.Name()])
				}
			}
		}
	}

	for k, v := range tps {
		g := GroupType{
			Group: k,
			Types: v,
		}
		groups = append(groups, g)
	}

	return groups
}
