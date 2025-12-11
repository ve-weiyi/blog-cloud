package ts

import (
	"sort"
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
	"github.com/ve-weiyi/ve-blog-golang/tools/parserx/apiparser/aspec"
)

func ConvertApiService(sp *aspec.ApiSpec) TsApiService {
	mgs := make(map[string][]TsApiGroup)
	for _, v := range sp.Service.Groups {
		name := v.Annotation.Properties["group"]
		if name == "" {
			name = sp.Service.Name
		}
		mgs[name] = append(mgs[name], ConvertApiGroup(v))
	}

	mts := make(map[string]TsType)
	for _, v := range sp.Types {
		mts[v.Name()] = convertTypeTs(v)
	}

	return TsApiService{
		ServiceName: sp.Service.Name,
		Groups:      mgs,
		Types:       mts,
	}
}

func ConvertApiGroup(g aspec.Group) TsApiGroup {
	routes := make([]TsApiRoute, 0, len(g.Routes))
	for _, r := range g.Routes {
		routes = append(routes, TsApiRoute{
			Summery:     buildDoc(r),
			Handler:     jsonconv.FirstLower(r.Handler),
			Path:        r.Path,
			Method:      strings.ToUpper(r.Method),
			Request:     getTypeName(r.RequestType),
			Response:    getTypeName(r.ResponseType),
			PathFields:  extractFields(r.PathType),
			QueryFields: extractFields(r.QueryType),
			FormFields:  extractFields(r.FormType),
		})
	}

	sort.Slice(routes, func(i, j int) bool {
		return routes[i].Path < routes[j].Path
	})

	return TsApiGroup{
		Prefix:     g.Annotation.Properties["prefix"],
		Middleware: parseMiddleware(g.Annotation.Properties["middleware"]),
		Routes:     routes,
	}
}

func buildDoc(r aspec.Route) string {
	var sb strings.Builder
	for _, d := range r.Doc {
		sb.WriteString(d)
	}
	for _, d := range r.HandlerDoc {
		sb.WriteString(d)
	}
	if r.AtDoc.Text != "" {
		sb.WriteString(strings.Trim(strings.Trim(r.AtDoc.Text, "\\"), "\""))
	}
	return sb.String()
}

func getTypeName(t aspec.Type) string {
	if t == nil {
		return "any"
	}
	return ConvertGoTypeToTsType(t.Name())
}

func extractFields(fields []aspec.Member) []string {
	if len(fields) == 0 {
		return nil
	}
	result := make([]string, 0, len(fields))
	for _, f := range fields {
		result = append(result, jsonconv.Case2Snake(f.Name))
	}
	return result
}

func parseMiddleware(middleware string) []string {
	if middleware == "" {
		return []string{}
	}
	return strings.Split(middleware, ",")
}

func convertTypeTs(st aspec.Type) TsType {
	t, ok := st.(aspec.DefineStruct)
	if !ok {
		return TsType{}
	}

	var extends []string
	fields := make([]TsTypeField, 0, len(t.Members))

	for _, v := range t.Members {
		if v.Name == "" {
			extends = append(extends, v.Type.Name())
		} else {
			fields = append(fields, TsTypeField{
				Comment:  v.Comment,
				Name:     jsonconv.Case2Snake(v.Name),
				Type:     ConvertGoTypeToTsType(v.Type.Name()),
				Nullable: strings.HasPrefix(v.Type.Name(), "*") || strings.Contains(v.Tag, "optional"),
			})
		}
	}

	return TsType{
		Comment: strings.Join(t.Comments(), "\n"),
		Name:    t.Name(),
		Extends: extends,
		Fields:  fields,
	}
}

// 转换go类型为ts类型
func ConvertGoTypeToTsType(name string) string {
	if strings.Contains(name, "map") {
		return "Record<string, any>"
	}
	if strings.HasPrefix(name, "*") {
		return ConvertGoTypeToTsType(name[1:]) // 指针
	}
	if strings.HasPrefix(name, "[]") {
		return ConvertGoTypeToTsType(name[2:]) + "[]" // 数组
	}
	if strings.LastIndex(name, ".") > 0 {
		return ConvertGoTypeToTsType(name[strings.LastIndex(name, ".")+1:]) // 去掉包名
	}
	switch name {
	case "int", "int32", "int64", "uint", "uint32", "uint64", "float32", "float64":
		return "number"
	case "string":
		return "string"
	case "bool":
		return "boolean"
	case "file":
		return "File"
	case "Time":
		return "string"
	case "FileHeader":
		return "File"
	case "interface{}", "object":
		return "any"
	default:
		return name
	}
}
