/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package web

import (
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ve-weiyi/ve-blog-golang/kit/tools/invent"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
	"github.com/ve-weiyi/ve-blog-golang/quickstart/cmd/web/helper"
	"github.com/ve-weiyi/ve-blog-golang/quickstart/tools/parserx"
	"github.com/ve-weiyi/ve-blog-golang/quickstart/tools/parserx/aspec"
)

// typescriptCmd represents the typescript command
var typescriptCmd = &cobra.Command{
	Use:   "typescript",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("typescript called")
		RunTypescriptCmd(cmd, args)
	},
}

func init() {
	ParseFlagVar(typescriptCmd)
}

func RunTypescriptCmd(cmd *cobra.Command, args []string) {
	RunTypescript(cmdVar)
}

func RunTypescript(conf *CmdVar) {
	log.Println("typescript called", jsonconv.AnyToJsonIndent(conf))
	var err error
	var sp *aspec.ApiSpec
	switch conf.VarStringMode {
	case "api":
		sp, err = parserx.NewSpecParser().ParseApi(conf.VarStringApiFile)
	case "swagger":
		sp, err = parserx.NewSwaggerParser().ParseApi(conf.VarStringApiFile)
	case "ast":
	}

	if err != nil {
		panic(err)
	}

	err = generateApiTs(sp, conf)
	if err != nil {
		panic(err)
	}

	err = generateTypesTs(sp, conf)
	if err != nil {
		panic(err)
	}
}

func generateApiTs(sp *aspec.ApiSpec, conf *CmdVar) error {
	t := path.Join(conf.VarStringTplPath, "api.ts.tpl")

	tpl, err := os.ReadFile(t)
	if err != nil {
		return err
	}

	groups := helper.ConvertApiTs(sp)

	var metas []invent.TemplateMeta
	for k, v := range groups {

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
			Key:            "",
			Mode:           invent.ModeCreateOrReplace,
			CodeOutPath:    path.Join(conf.VarStringOutPath, fmt.Sprintf("%s.ts", k)),
			TemplateString: string(tpl),
			FunMap: map[string]any{
				"Join": func(s []string) string {
					return strings.Join(s, ", ")
				},
			},
			Data: map[string]any{
				"ImportPkgPaths": []string{
					`import request from "@/utils/request";`,
					fmt.Sprintf(`import { %s } from "./types";`, strings.Join(ims, ", ")),
				},
				"GroupRoutes": v,
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

func generateTypesTs(sp *aspec.ApiSpec, conf *CmdVar) error {
	t := path.Join(conf.VarStringTplPath, "types.ts.tpl")
	o := path.Join(conf.VarStringOutPath, "types.ts")

	tpl, err := os.ReadFile(t)
	if err != nil {
		return err
	}

	//ims := strings.Split(conf.VarStringIgnoreModel, ",")
	ts := helper.ConvertTypeTs(sp)

	meta := invent.TemplateMeta{
		Mode:           invent.ModeCreateOrReplace,
		CodeOutPath:    o,
		TemplateString: string(tpl),
		FunMap: map[string]any{
			"Join": func(s []string) string {
				return strings.Join(s, ", ")
			},
		},
		Data: ts,
	}

	return meta.Execute()
}
