package apiparser

import (
	"os"
	"path"

	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"

	"github.com/ve-weiyi/ve-blog-golang/tools/parserx/apiparser/aspec"
)

type SpecParser struct {
}

func NewSpecParser() ApiParser {
	return &SpecParser{}
}

func (s *SpecParser) ParseApi(filename string) (out *aspec.ApiSpec, err error) {
	sp, err := ParseApiSpec(filename)
	if err != nil {
		return nil, err
	}

	gps := getGroupsFromSpec(sp)
	tps := getTypesFromSpec(sp)

	out = &aspec.ApiSpec{
		Service: aspec.Service{
			Name:   sp.Service.Name,
			Groups: gps,
		},
		Types: tps,
	}

	return out, nil
}

func ParseApiSpec(filename string) (out *spec.ApiSpec, err error) {
	if path.IsAbs(filename) {
		return parser.Parse(filename)
	}

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	f := path.Join(dir, filename)
	return parser.Parse(f)
}

func getGroupsFromSpec(sp *spec.ApiSpec) []aspec.Group {
	var gps []aspec.Group
	for _, g := range sp.Service.Groups {
		var rts []aspec.Route
		//var base = g.Annotation.Properties["prefix"]
		for _, r := range g.Routes {

			rt := aspec.Route{
				Method:         r.Method,
				Path:           r.Path,
				PathType:       nil,
				HeaderType:     nil,
				FormType:       nil,
				QueryType:      nil,
				RequestType:    r.RequestType,
				ResponseType:   r.ResponseType,
				Docs:           nil,
				Handler:        r.Handler,
				AtDoc:          aspec.AtDoc(r.AtDoc),
				HandlerDoc:     nil,
				HandlerComment: nil,
				Doc:            aspec.Doc(r.Doc),
				Comment:        nil,
			}

			rts = append(rts, rt)
		}

		gp := aspec.Group{
			Annotation: aspec.Annotation(g.Annotation),
			Routes:     rts,
		}

		gps = append(gps, gp)
	}
	return gps
}

func getTypesFromSpec(sp *spec.ApiSpec) []aspec.Type {
	var types []aspec.Type
	for _, st := range sp.Types {
		var tp aspec.Type
		tp = convertType(st)
		types = append(types, tp)
	}

	return types
}

func convertType(st spec.Type) aspec.Type {
	var out aspec.Type
	switch t := st.(type) {
	case spec.DefineStruct:

		var ms []aspec.Member
		for _, member := range t.Members {
			ms = append(ms, aspec.Member{
				Name:     member.Name,
				Type:     convertType(member.Type),
				Tag:      member.Tag,
				Comment:  member.Comment,
				Docs:     aspec.Doc(member.Docs),
				IsInline: member.IsInline,
			})
		}

		out = aspec.DefineStruct{
			RawName: t.RawName,
			Members: ms,
			Docs:    aspec.Doc(t.Docs),
		}
	case spec.PrimitiveType:
		out = aspec.PrimitiveType{
			RawName: t.RawName,
		}
	case spec.MapType:
		out = aspec.MapType{
			RawName: t.RawName,
			Key:     t.Key,
			Value:   convertType(t.Value),
		}
	case spec.ArrayType:
		out = aspec.ArrayType{
			RawName: t.RawName,
			Value:   convertType(t.Value),
		}
	case spec.InterfaceType:
		out = aspec.InterfaceType{
			RawName: t.RawName,
		}
	case spec.PointerType:
		out = aspec.PointerType{
			RawName: t.RawName,
			Type:    convertType(t.Type),
		}
	default:
		panic("unknown type")
	}

	return out
}
