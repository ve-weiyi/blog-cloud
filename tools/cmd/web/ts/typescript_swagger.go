package ts

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ve-weiyi/ve-blog-golang/tools/parserx/apiparser"
)

var typescriptSwaggerFlags = struct {
	SqlFile string // sql文件
	ApiFile string // api文件

	TplPath     string // 模板路径
	OutPath     string // 输出路径
	NameAs      string // 输出名称
	Mode        string // 解析模式 swagger、api、ast
	IgnoreModel string // 忽略的模型
}{
	SqlFile: "vars.sql",
	ApiFile: "test.api",
	TplPath: "test.tpl",
	OutPath: "./",
	NameAs:  "%s.go",
	Mode:    "api",
}

func NewTypescriptSwaggerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swagger",
		Short: "根据 api 文件生成 TypeScript 代码",
		RunE:  runTypescriptSwagger,
	}

	cmd.Flags().StringVarP(&typescriptSwaggerFlags.ApiFile, "api-file", "f", typescriptSwaggerFlags.ApiFile, "API文件路径")
	cmd.Flags().StringVarP(&typescriptSwaggerFlags.TplPath, "tpl-path", "t", typescriptSwaggerFlags.TplPath, "模板文件路径")
	cmd.Flags().StringVarP(&typescriptSwaggerFlags.OutPath, "out-path", "o", typescriptSwaggerFlags.OutPath, "输出目录路径")
	cmd.Flags().StringVarP(&typescriptSwaggerFlags.NameAs, "name-as", "n", typescriptSwaggerFlags.NameAs, "输出文件名称")
	cmd.Flags().StringVarP(&typescriptSwaggerFlags.Mode, "mode", "m", typescriptSwaggerFlags.Mode, "解析模式")
	cmd.Flags().StringVarP(&typescriptSwaggerFlags.IgnoreModel, "ignore-model", "i", typescriptSwaggerFlags.IgnoreModel, "忽略的模型")

	return cmd
}

func runTypescriptSwagger(cmd *cobra.Command, args []string) error {
	fmt.Println("===== 命令参数 =====")
	fmt.Printf("api-file: %s\n", typescriptSwaggerFlags.ApiFile)
	fmt.Printf("tpl-path: %s\n", typescriptSwaggerFlags.TplPath)
	fmt.Printf("out-path: %s\n", typescriptSwaggerFlags.OutPath)
	fmt.Printf("name-as: %s\n", typescriptSwaggerFlags.NameAs)
	fmt.Println("====================")

	sp, err := apiparser.NewSwaggerParser().ParseApi(typescriptSwaggerFlags.ApiFile)
	if err != nil {
		return err
	}

	sv := ConvertApiService(sp)

	if err = generateApiTs(sv, typescriptSwaggerFlags.TplPath, typescriptSwaggerFlags.OutPath, typescriptSwaggerFlags.NameAs); err != nil {
		return err
	}

	if err = generateTypesTs(sv, typescriptSwaggerFlags.TplPath, typescriptSwaggerFlags.OutPath, typescriptSwaggerFlags.NameAs); err != nil {
		return err
	}

	fmt.Println("TypeScript code generated successfully")
	return nil
}
