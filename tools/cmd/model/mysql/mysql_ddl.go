package mysql

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ve-weiyi/ve-blog-golang/tools/parserx/dbparser"
)

var ddlFlags = struct {
	SqlFile string
	TplFile string
	OutPath string
	NameAs  string
}{
	SqlFile: "test.sql",
	TplFile: "model.tpl",
	OutPath: "./tem",
	NameAs:  "%s.go",
}

func NewMysqlDDLCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ddl",
		Short: "从 sql 文件生成 model 代码",
		RunE:  runMysqlDDL,
	}

	cmd.Flags().StringVarP(&ddlFlags.SqlFile, "sql-file", "s", ddlFlags.SqlFile, "sql文件路径")
	cmd.Flags().StringVarP(&ddlFlags.TplFile, "tpl-file", "t", ddlFlags.TplFile, "模板文件路径")
	cmd.Flags().StringVarP(&ddlFlags.OutPath, "out-path", "o", ddlFlags.OutPath, "输出目录路径")
	cmd.Flags().StringVarP(&ddlFlags.NameAs, "name-as", "n", ddlFlags.NameAs, "输出文件名称")

	return cmd
}

func runMysqlDDL(cmd *cobra.Command, args []string) error {
	// 打印获取到的 flag 参数
	fmt.Println("===== 命令参数 =====")
	fmt.Printf("sql-file: %s\n", ddlFlags.SqlFile)
	fmt.Printf("tpl-file: %s\n", ddlFlags.TplFile)
	fmt.Printf("out-path: %s\n", ddlFlags.OutPath)
	fmt.Printf("name-as: %s\n", ddlFlags.NameAs)
	fmt.Println("====================")

	tables, err := dbparser.ParseTableFromSql(ddlFlags.SqlFile)
	if err != nil {
		return err
	}

	var models []*ModelData
	for _, table := range tables {
		m := ConvertTableToData(table)
		models = append(models, m)
	}

	err = generateModel(models, ddlFlags.TplFile, ddlFlags.OutPath, ddlFlags.NameAs)
	if err != nil {
		return err
	}

	fmt.Println("Model code generated successfully")
	return nil
}
