package gin

import (
	"sort"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/util"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
	"github.com/ve-weiyi/ve-blog-golang/tools/parserx/apiparser/aspec"
)

func ConvertRouteGroups(sp *aspec.ApiSpec) (out map[string][]GroupRoute) {
	var defaultGroupName = sp.Service.Name
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
			name = defaultGroupName
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

// 简化版的 groupTypes 生成逻辑
// 参考go-zero的type分组实现
func GroupTypes(api *aspec.ApiSpec) map[string]map[string]aspec.Type {
	groupTypeDefault := "types"
	// 步骤 1: 收集每个类型被哪些分组引用
	typeToGroups := make(map[string]map[string]bool) // key: typeName, value: map[groupName]bool

	for _, group := range api.Service.Groups {
		// 获取分组名
		groupName := group.GetAnnotation("group")
		if groupName == "" {
			groupName = groupTypeDefault
		}
		groupName = util.SafeString(strings.TrimSuffix(strings.TrimPrefix(groupName, "/"), "/"))

		for _, route := range group.Routes {
			// 处理 RequestType
			if route.RequestType != nil {
				typeName := getTypeName(route.RequestType)
				if typeName != "" {
					if _, ok := typeToGroups[typeName]; !ok {
						typeToGroups[typeName] = make(map[string]bool)
					}
					typeToGroups[typeName][groupName] = true
				}
			}
			// 处理 ResponseType
			if route.ResponseType != nil {
				typeName := getTypeName(route.ResponseType)
				if typeName != "" {
					if _, ok := typeToGroups[typeName]; !ok {
						typeToGroups[typeName] = make(map[string]bool)
					}
					typeToGroups[typeName][groupName] = true
				}
			}
		}
	}

	// 步骤 2: 构建并填充 groupTypes
	groupTypes := make(map[string]map[string]aspec.Type)

	// 初始化默认分组
	groupTypes[groupTypeDefault] = make(map[string]aspec.Type)

	for _, typ := range api.Types {
		typeName := util.Title(typ.Name())

		// 获取引用该类型的所有分组
		groups := typeToGroups[typeName]
		typeCount := len(groups)

		var targetGroup string

		if typeCount == 1 {
			// 情况 B: 只被一个分组引用，就归属于该分组
			// 从 map 中取出唯一的 key
			for g := range groups {
				targetGroup = g
				break
			}
			// 确保目标分组已初始化
			if _, ok := groupTypes[targetGroup]; !ok {
				groupTypes[targetGroup] = make(map[string]aspec.Type)
			}
		} else {
			// 情况 A (typeCount == 0) 或 情况 C (typeCount > 1): 都归为默认分组
			targetGroup = groupTypeDefault
		}

		// 将类型放入目标分组
		groupTypes[targetGroup][typeName] = typ
	}

	// 正确的方式：使用索引遍历
	for groupName := range groupTypes {
		// 通过键获取原始的 map
		typesMap := groupTypes[groupName]

		// 1. 提取所有类型名并排序
		typeNames := make([]string, 0, len(typesMap))
		for name := range typesMap {
			typeNames = append(typeNames, name)
		}
		sort.Strings(typeNames)

		// 2. 创建一个新的、有序的 map
		sortedTypesMap := make(map[string]aspec.Type, len(typesMap))
		for _, name := range typeNames {
			sortedTypesMap[name] = typesMap[name]
		}

		// 3. 将排序后的 map 重新赋值回 groupTypes
		groupTypes[groupName] = sortedTypesMap
	}

	return groupTypes
}

// getTypeName 是一个辅助函数，用于从 aspec.Type 中获取简化的类型名
// 这是一个简化版，你可能需要根据实际的 aspec.Type 结构来实现
func getTypeName(t aspec.Type) string {
	if t == nil {
		return ""
	}
	// 假设 t.Name() 可以返回类型的名称，如 "CreateUserRequest"
	// 在实际的 go-zero 代码中，逻辑可能更复杂，因为 t 可能是一个引用类型
	return util.Title(t.Name())
}
