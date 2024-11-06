package web

import (
	"github.com/spf13/cobra"
)

type CmdVar struct {
	VarStringSqlFile string // sql文件
	VarStringApiFile string // api文件

	VarStringTplPath     string // 模板路径
	VarStringOutPath     string // 输出路径
	VarStringNameAs      string // 输出名称
	VarStringMode        string // 解析模式 swagger、api、ast
	VarStringIgnoreModel string // 忽略的模型
}

var cmdVar = &CmdVar{}

func ParseFlagVar(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&cmdVar.VarStringApiFile, "api-file", "f", "test.api", "sql文件")
	cmd.PersistentFlags().StringVarP(&cmdVar.VarStringTplPath, "tpl-path", "t", "test.tpl", "模板文件")
	cmd.PersistentFlags().StringVarP(&cmdVar.VarStringOutPath, "out-path", "o", "./", "输出路径")
	cmd.PersistentFlags().StringVarP(&cmdVar.VarStringNameAs, "name-as", "n", "%s.go", "输出名称")
	cmd.PersistentFlags().StringVarP(&cmdVar.VarStringMode, "mode", "m", "api", "解析模式")
	cmd.PersistentFlags().StringVarP(&cmdVar.VarStringIgnoreModel, "ignore-model", "i", "", "忽略的模型")
}
