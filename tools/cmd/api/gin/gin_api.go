package gin

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ve-weiyi/ve-blog-golang/tools/parserx/apiparser"
)

var apiFlags = struct {
	ApiFile        string
	TplPath        string
	OutPath        string
	NameAs         string
	ContextPackage string
}{
	ApiFile:        "test.api",
	TplPath:        "test.tpl",
	OutPath:        "./",
	NameAs:         "%s.go",
	ContextPackage: "context",
}

func NewGinApiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "从 .api文件生成 Gin 框架代码",
		RunE:  runApi,
	}

	cmd.Flags().StringVarP(&apiFlags.ApiFile, "api-file", "f", apiFlags.ApiFile, "API文件路径")
	cmd.Flags().StringVarP(&apiFlags.TplPath, "tpl-path", "t", apiFlags.TplPath, "模板路径")
	cmd.Flags().StringVarP(&apiFlags.OutPath, "out-path", "o", apiFlags.OutPath, "输出路径")
	cmd.Flags().StringVarP(&apiFlags.NameAs, "name-as", "n", apiFlags.NameAs, "输出名称")
	cmd.Flags().StringVarP(&apiFlags.ContextPackage, "svctx-package", "c", swagFlags.ContextPackage, "导入上下文包名")

	return cmd
}

func runApi(cmd *cobra.Command, args []string) error {
	fmt.Println("===== 命令参数 =====")
	fmt.Printf("api-file: %s\n", apiFlags.ApiFile)
	fmt.Printf("tpl-path: %s\n", apiFlags.TplPath)
	fmt.Printf("out-path: %s\n", apiFlags.OutPath)
	fmt.Printf("name-as: %s\n", apiFlags.NameAs)
	fmt.Printf("context-package: %s\n", apiFlags.ContextPackage)
	fmt.Println("====================")

	sp, err := apiparser.NewSpecParser().ParseApi(apiFlags.ApiFile)
	if err != nil {
		return err
	}

	err = generateTypes(sp, apiFlags.TplPath, apiFlags.OutPath, apiFlags.NameAs)
	if err != nil {
		return err
	}

	err = generateLogics(sp, apiFlags.TplPath, apiFlags.OutPath, apiFlags.NameAs, apiFlags.ContextPackage)
	if err != nil {
		return err
	}

	err = generateHandlers(sp, apiFlags.TplPath, apiFlags.OutPath, apiFlags.NameAs, apiFlags.ContextPackage)
	if err != nil {
		return err
	}

	err = generateRouters(sp, apiFlags.TplPath, apiFlags.OutPath, apiFlags.NameAs, apiFlags.ContextPackage)
	if err != nil {
		return err
	}

	err = generateRoutes(sp, apiFlags.TplPath, apiFlags.OutPath, apiFlags.NameAs, apiFlags.ContextPackage)
	if err != nil {
		return err
	}

	fmt.Println("Gin code generated successfully")
	return nil
}
