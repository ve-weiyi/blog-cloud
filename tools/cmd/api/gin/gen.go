package gin

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"

	"github.com/ve-weiyi/ve-blog-golang/pkg/kit/quickstart/gotplgen"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/jsonconv"
	"github.com/ve-weiyi/ve-blog-golang/tools/parserx/apiparser/aspec"
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

type GroupType map[string]map[string]aspec.Type

func generateHandlers(sp *aspec.ApiSpec, tplPath, outPath, nameAs, contextPackage string) error {
	var metas []gotplgen.TemplateMeta

	tpl, err := os.ReadFile(path.Join(tplPath, "handler.tpl"))
	if err != nil {
		return err
	}

	pkg, _, _ := golang.GetParentPackage(outPath)
	var groups map[string][]GroupRoute
	groups = ConvertRouteGroups(sp)

	for k, v := range groups {

		m := gotplgen.TemplateMeta{
			Mode:           gotplgen.ModeCreateOrReplace,
			CodeOutPath:    path.Join(outPath, "handler", fmt.Sprintf("%v_handler.go", k)),
			TemplateString: string(tpl),
			Data: map[string]any{
				"Package": "handler",
				"Imports": []string{
					fmt.Sprintf(`"%s"`, contextPackage),
					fmt.Sprintf(`"%s/types"`, pkg),
					fmt.Sprintf(`"%s/logic"`, pkg),
				},
				"Group":       jsonconv.Case2Camel(k),
				"GroupRoutes": v,
			},
			FunMap: map[string]any{
				"pkgTypes":     pkgTypesFunc,
				"commentTypes": commentTypesFunc,
			},
		}

		metas = append(metas, m)
	}

	for _, m := range metas {
		err := m.Execute()
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func generateLogics(sp *aspec.ApiSpec, tplPath, outPath, nameAs, contextPackage string) error {
	var metas []gotplgen.TemplateMeta

	tpl, err := os.ReadFile(path.Join(tplPath, "logic.tpl"))
	if err != nil {
		return err
	}

	pkg, _, _ := golang.GetParentPackage(outPath)

	var groups map[string][]GroupRoute
	groups = ConvertRouteGroups(sp)

	for k, v := range groups {

		m := gotplgen.TemplateMeta{
			Mode:           gotplgen.ModeCreateOrReplace,
			CodeOutPath:    path.Join(outPath, "logic", fmt.Sprintf("%v_logic.go", k)),
			TemplateString: string(tpl),
			Data: map[string]any{
				"Package": "logic",
				"Imports": []string{
					fmt.Sprintf(`"%s"`, contextPackage),
					fmt.Sprintf(`"%s/types"`, pkg),
				},
				"Group":       jsonconv.Case2Camel(k),
				"GroupRoutes": v,
			},
			FunMap: map[string]any{
				"pkgTypes": pkgTypesFunc,
			},
		}

		metas = append(metas, m)
	}
	for _, m := range metas {
		err := m.Execute()
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func generateRouters(sp *aspec.ApiSpec, tplPath, outPath, nameAs, contextPackage string) error {
	var metas []gotplgen.TemplateMeta

	tpl, err := os.ReadFile(path.Join(tplPath, "router.tpl"))
	if err != nil {
		return err
	}

	pkg, _, _ := golang.GetParentPackage(outPath)

	var groups map[string][]GroupRoute
	groups = ConvertRouteGroups(sp)

	for k, v := range groups {

		m := gotplgen.TemplateMeta{
			Mode:           gotplgen.ModeCreateOrReplace,
			CodeOutPath:    path.Join(outPath, "router", fmt.Sprintf("%v_router.go", k)),
			TemplateString: string(tpl),
			Data: map[string]any{
				"Package": "router",
				"Imports": []string{
					fmt.Sprintf(`"%s"`, contextPackage),
					fmt.Sprintf(`"%s/types"`, pkg),
					fmt.Sprintf(`"%s/handler"`, pkg),
				},
				"Group":       jsonconv.Case2Camel(k),
				"GroupRoutes": v,
			},
		}

		metas = append(metas, m)
	}

	for _, m := range metas {
		err := m.Execute()
		if err != nil {
			return err
		}
	}

	return nil
}

func generateRoutes(sp *aspec.ApiSpec, tplPath, outPath, nameAs, contextPackage string) error {
	var metas []gotplgen.TemplateMeta

	tpl, err := os.ReadFile(path.Join(tplPath, "routes.tpl"))
	if err != nil {
		return err
	}

	pkg, _, _ := golang.GetParentPackage(outPath)

	var groups map[string][]GroupRoute
	groups = ConvertRouteGroups(sp)

	var gps []string
	for k, _ := range groups {
		gps = append(gps, jsonconv.Case2Camel(k))
	}

	sort.Slice(gps, func(i, j int) bool {
		return gps[i] < gps[j]
	})

	m := gotplgen.TemplateMeta{
		Mode:           gotplgen.ModeCreateOrReplace,
		CodeOutPath:    path.Join(outPath, "routes.go"),
		TemplateString: string(tpl),
		Data: map[string]any{
			"Package": filepath.Base(outPath),
			"Imports": []string{
				fmt.Sprintf(`"%s"`, contextPackage),
				fmt.Sprintf(`"%s/router"`, pkg),
			},
			"Groups": gps,
		},
		FunMap: gotplgen.StdMapUtils,
	}

	metas = append(metas, m)

	for _, m := range metas {
		err := m.Execute()
		if err != nil {
			return err
		}
	}

	return nil
}

func generateTypes(sp *aspec.ApiSpec, tplPath, outPath, nameAs string) error {
	handlerTpl, err := os.ReadFile(path.Join(tplPath, "types.tpl"))
	if err != nil {
		return err
	}

	groupTypes := GroupTypes(sp)

	var metas []gotplgen.TemplateMeta
	for k, typeGroup := range groupTypes {
		if len(typeGroup) == 0 {
			continue
		}

		ts := make([]aspec.Type, 0, len(typeGroup))
		for _, v := range typeGroup {
			ts = append(ts, v)
		}

		sort.Slice(ts, func(i, j int) bool {
			return ts[i].Name() < ts[j].Name()
		})

		var types []string
		for _, t := range ts {
			val, err := buildTypes(t)
			if err != nil {
				return err
			}

			types = append(types, val)
		}

		meta := gotplgen.TemplateMeta{
			Mode:           gotplgen.ModeCreateOrReplace,
			CodeOutPath:    path.Join(outPath, "types", fmt.Sprintf("%v.go", k)),
			TemplateString: string(handlerTpl),
			Data: map[string]any{
				"Package": "types",
				"Imports": []string{},
				"Name":    jsonconv.Case2Camel(k),
				"Types":   types,
			},
		}

		metas = append(metas, meta)
	}

	for _, m := range metas {
		err := m.Execute()
		if err != nil {
			return err
		}
	}
	return nil
}

func buildTypes(tp aspec.Type) (string, error) {
	var builder strings.Builder

	structType, ok := tp.(aspec.DefineStruct)
	if !ok {
		return "", fmt.Errorf("unspport struct type: %s", tp.Name())
	}

	builder.WriteString(fmt.Sprintf("type %s struct {\n", jsonconv.Case2Camel(structType.RawName)))

	for _, member := range structType.Members {
		if member.IsInline {
			builder.WriteString(fmt.Sprintf("%v\n", member.Type.Name()))
			continue
		}

		builder.WriteString(fmt.Sprintf("	%s %s %s", member.Name, member.Type.Name(), member.Tag))

		var comment = member.GetComment()
		if len(comment) > 0 {
			comment = strings.TrimPrefix(comment, "//")
			comment = "//" + comment
		}

		if len(comment) > 0 {
			builder.WriteString(fmt.Sprintf("	%s", comment))
		}

		builder.WriteString("\n")
	}

	builder.WriteString("}")

	return builder.String(), nil
}

func pkgTypesFunc(input string) string {
	re := regexp.MustCompile(`\w+`)
	result := re.ReplaceAllString(input, "types.$0")

	if strings.HasPrefix(result, "[]") {
		return result
	}

	return fmt.Sprintf("*%v", result)
}

func commentTypesFunc(input string) string {
	re := regexp.MustCompile(`\w+`)
	result := re.ReplaceAllString(input, "types.$0")
	return strings.ReplaceAll(result, "*", "")
}
