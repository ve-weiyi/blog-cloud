package ts

import (
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/kit/quickstart/invent"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
)

type TsApiService struct {
	ServiceName string
	Groups      map[string][]TsApiGroup // [account][a,b,c]
	Types       map[string]TsType
}

type TsApiGroup struct {
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

	PathFields  []string
	QueryFields []string
	FormFields  []string
	BodyFields  []string
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

func generateApiTs(sv TsApiService, tplPath string, outPath string, nameAs string) error {
	t := path.Join(tplPath, "api.ts.tpl")

	tpl, err := os.ReadFile(t)
	if err != nil {
		return err
	}

	var metas []invent.TemplateMeta
	for k, v := range sv.Groups {

		mt := make(map[string]string)
		for _, g := range v {
			for _, r := range g.Routes {
				if r.Request != "" {
					name := strings.Trim(r.Request, "[]")
					mt[name] = name
				}
				if r.Response != "" {
					name := strings.Trim(r.Response, "[]")
					mt[name] = name
				}
			}
		}

		var ims []string
		for kk, _ := range mt {
			if kk == "any" {
				continue
			}
			ims = append(ims, kk)
		}

		sort.Slice(ims, func(i, j int) bool {
			return ims[i] < ims[j]
		})

		meta := invent.TemplateMeta{
			Mode:           invent.ModeCreateOrReplace,
			CodeOutPath:    path.Join(outPath, fmt.Sprintf("%s.ts", k)),
			TemplateString: string(tpl),
			FunMap: map[string]any{
				"Join": func(s []string) string {
					return strings.Join(s, ", ")
				},
			},
			Data: map[string]any{
				"ImportPkgPaths": []string{
					`import request from "@/utils/request";`,
					fmt.Sprintf(`import type { %s } from "./types";`, strings.Join(ims, ", ")),
				},
				"Name":        jsonconv.Case2Camel(k),
				"TsApiGroups": v,
			},
		}

		metas = append(metas, meta)
	}

	for _, m := range metas {
		err := m.Execute()
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func generateTypesTs(sv TsApiService, tplPath string, outPath string, nameAs string) error {
	t := path.Join(tplPath, "types.ts.tpl")
	o := path.Join(outPath, "types.ts")

	tpl, err := os.ReadFile(t)
	if err != nil {
		return err
	}

	meta := invent.TemplateMeta{
		Mode:           invent.ModeCreateOrReplace,
		CodeOutPath:    o,
		TemplateString: string(tpl),
		FunMap: map[string]any{
			"Join": func(s []string) string {
				return strings.Join(s, ", ")
			},
		},
		Data: sv.Types,
	}

	return meta.Execute()
}

func generatePermsTs(sv TsApiService, tplPath string, outPath string, nameAs string) error {
	t := path.Join(tplPath, "perms.ts.tpl")
	o := path.Join(outPath, "perms.ts")

	tpl, err := os.ReadFile(t)
	if err != nil {
		return err
	}

	meta := invent.TemplateMeta{
		Mode:           invent.ModeCreateOrReplace,
		CodeOutPath:    o,
		TemplateString: string(tpl),
		FunMap:         nil,
		Data: map[string]any{
			"Name":   jsonconv.Case2Camel(sv.ServiceName),
			"Groups": sv.Groups,
		},
	}

	return meta.Execute()
}
