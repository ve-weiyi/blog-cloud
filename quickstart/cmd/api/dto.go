/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package api

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/api/gogen"

	"github.com/ve-weiyi/ve-blog-golang/kit/tools/invent"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
	"github.com/ve-weiyi/ve-blog-golang/quickstart/cmd/api/helper"
	"github.com/ve-weiyi/ve-blog-golang/quickstart/tools/parserx"
)

// dtoCmd represents the dto command
var dtoCmd = &cobra.Command{
	Use:   "dto",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := generateDto(cmdVar)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	ParseFlagVar(dtoCmd)
}

func generateDto(cv *CmdVar) error {
	sp, err := parserx.ParseApiSpec(cv.VarStringApiFile)
	if err != nil {
		return err
	}

	var metas []invent.TemplateMeta

	handlerTpl, err := os.ReadFile(path.Join(cv.VarStringTplPath, "dto.tpl"))
	if err != nil {
		return err
	}

	var groups []helper.GroupType
	groups = append(groups, helper.GroupType{
		Group: "dto",
		Types: sp.Types,
	})

	for _, v := range groups {
		val, err := gogen.BuildTypes(v.Types)
		if err != nil {
			return err
		}

		meta := invent.TemplateMeta{
			Mode:           invent.ModeCreateOrReplace,
			CodeOutPath:    path.Join(cv.VarStringOutPath, "dto", fmt.Sprintf(cv.VarStringNameAs, v.Group)),
			TemplateString: string(handlerTpl),
			Data: map[string]any{
				"Package": "dto",
				"Imports": []string{},
				"Name":    jsonconv.Case2Camel(v.Group),
				"Types":   val,
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
