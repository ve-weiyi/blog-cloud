package gin

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ve-weiyi/ve-blog-golang/tools/parserx/apiparser"
)

var swagFlags = struct {
	SwaggerFile    string
	TplPath        string
	OutPath        string
	NameAs         string
	ContextPackage string
}{
	SwaggerFile:    "test.api",
	TplPath:        "test.tpl",
	OutPath:        "./",
	NameAs:         "%s.go",
	ContextPackage: "context",
}

func NewGinSwaggerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swagger",
		Short: "从 swagger.json文件生成 Gin 框架代码",
		RunE:  runSwagger,
	}

	cmd.Flags().StringVarP(&swagFlags.SwaggerFile, "api-file", "f", swagFlags.SwaggerFile, "API文件路径")
	cmd.Flags().StringVarP(&swagFlags.TplPath, "tpl-path", "t", swagFlags.TplPath, "模板路径")
	cmd.Flags().StringVarP(&swagFlags.OutPath, "out-path", "o", swagFlags.OutPath, "输出路径")
	cmd.Flags().StringVarP(&swagFlags.NameAs, "name-as", "n", swagFlags.NameAs, "输出名称")
	cmd.Flags().StringVarP(&swagFlags.ContextPackage, "svctx-package", "c", swagFlags.ContextPackage, "导入上下文包名")

	return cmd
}

func runSwagger(cmd *cobra.Command, args []string) error {
	fmt.Println("===== 命令参数 =====")
	fmt.Printf("swagger-file: %s\n", swagFlags.SwaggerFile)
	fmt.Printf("tpl-path: %s\n", swagFlags.TplPath)
	fmt.Printf("out-path: %s\n", swagFlags.OutPath)
	fmt.Printf("name-as: %s\n", swagFlags.NameAs)
	fmt.Printf("context-package: %s\n", swagFlags.ContextPackage)
	fmt.Println("====================")

	sp, err := apiparser.NewSwaggerParser().ParseApi(swagFlags.SwaggerFile)
	if err != nil {
		return err
	}

	err = generateTypes(sp, swagFlags.TplPath, swagFlags.OutPath, swagFlags.NameAs)
	if err != nil {
		return err
	}

	err = generateLogics(sp, swagFlags.TplPath, swagFlags.OutPath, swagFlags.NameAs, swagFlags.ContextPackage)
	if err != nil {
		return err
	}

	err = generateHandlers(sp, swagFlags.TplPath, swagFlags.OutPath, swagFlags.NameAs, swagFlags.ContextPackage)
	if err != nil {
		return err
	}

	err = generateRouters(sp, swagFlags.TplPath, swagFlags.OutPath, swagFlags.NameAs, swagFlags.ContextPackage)
	if err != nil {
		return err
	}

	err = generateRoutes(sp, swagFlags.TplPath, swagFlags.OutPath, swagFlags.NameAs, swagFlags.ContextPackage)
	if err != nil {
		return err
	}

	fmt.Println("Gin code generated successfully")
	return nil
}
