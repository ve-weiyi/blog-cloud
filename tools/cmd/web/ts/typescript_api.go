package ts

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ve-weiyi/ve-blog-golang/tools/parserx/apiparser"
)

var typescriptApiFlags = struct {
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

func NewTypescriptApiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "根据 .api 文件生成 TypeScript 代码",
		RunE:  runTypescriptApi,
	}

	cmd.Flags().StringVarP(&typescriptApiFlags.ApiFile, "api-file", "f", typescriptApiFlags.ApiFile, "API文件路径")
	cmd.Flags().StringVarP(&typescriptApiFlags.TplPath, "tpl-path", "t", typescriptApiFlags.TplPath, "模板文件路径")
	cmd.Flags().StringVarP(&typescriptApiFlags.OutPath, "out-path", "o", typescriptApiFlags.OutPath, "输出目录路径")
	cmd.Flags().StringVarP(&typescriptApiFlags.NameAs, "name-as", "n", typescriptApiFlags.NameAs, "输出文件名称")
	cmd.Flags().StringVarP(&typescriptApiFlags.Mode, "mode", "m", typescriptApiFlags.Mode, "解析模式")
	cmd.Flags().StringVarP(&typescriptApiFlags.IgnoreModel, "ignore-model", "i", typescriptApiFlags.IgnoreModel, "忽略的模型")

	return cmd
}

func runTypescriptApi(cmd *cobra.Command, args []string) error {
	fmt.Println("===== 命令参数 =====")
	fmt.Printf("api-file: %s\n", typescriptApiFlags.ApiFile)
	fmt.Printf("tpl-path: %s\n", typescriptApiFlags.TplPath)
	fmt.Printf("out-path: %s\n", typescriptApiFlags.OutPath)
	fmt.Printf("name-as: %s\n", typescriptApiFlags.NameAs)
	fmt.Println("====================")

	sp, err := apiparser.NewSpecParser().ParseApi(typescriptApiFlags.ApiFile)
	if err != nil {
		panic(err)
	}

	sv := ConvertApiService(sp)

	if err = generateApiTs(sv, typescriptApiFlags.TplPath, typescriptApiFlags.OutPath, typescriptApiFlags.NameAs); err != nil {
		return err
	}

	if err = generateTypesTs(sv, typescriptApiFlags.TplPath, typescriptApiFlags.OutPath, typescriptApiFlags.NameAs); err != nil {
		return err
	}

	fmt.Println("TypeScript code generated successfully")
	return nil
}
