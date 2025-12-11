package mysql

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ve-weiyi/ve-blog-golang/tools/parserx/dbparser"
)

var dsnFlags = struct {
	Url     string
	TplFile string
	OutPath string
	NameAs  string
}{
	Url:     "root:123456@(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
	TplFile: "model.tpl",
	OutPath: "./tem",
	NameAs:  "%s.go",
}

func NewMysqlDSNCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dsn",
		Short: "从数据库连接生成 model 代码",
		RunE:  runMysqlDSN,
	}

	cmd.Flags().StringVarP(&dsnFlags.Url, "url", "u", dsnFlags.Url, "数据库地址")
	cmd.Flags().StringVarP(&dsnFlags.TplFile, "tpl-file", "t", dsnFlags.TplFile, "模板文件路径")
	cmd.Flags().StringVarP(&dsnFlags.OutPath, "out-path", "o", dsnFlags.OutPath, "输出目录路径")
	cmd.Flags().StringVarP(&dsnFlags.NameAs, "name-as", "n", dsnFlags.NameAs, "输出文件名称")

	return cmd
}

func runMysqlDSN(cmd *cobra.Command, args []string) error {
	fmt.Println("===== 命令参数 =====")
	fmt.Printf("url: %s\n", dsnFlags.Url)
	fmt.Printf("tpl-file: %s\n", dsnFlags.TplFile)
	fmt.Printf("out-path: %s\n", dsnFlags.OutPath)
	fmt.Printf("name-as: %s\n", dsnFlags.NameAs)
	fmt.Println("====================")

	tables, err := dbparser.ParseTableFromDsn(dsnFlags.Url)
	if err != nil {
		return err
	}

	var models []*ModelData
	for _, table := range tables {
		m := ConvertTableToData(table)
		models = append(models, m)
	}

	err = generateModel(models, dsnFlags.TplFile, dsnFlags.OutPath, dsnFlags.NameAs)
	if err != nil {
		return err
	}

	fmt.Println("Model code generated successfully")
	return nil
}
