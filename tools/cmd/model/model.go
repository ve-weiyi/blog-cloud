package model

import (
	"github.com/spf13/cobra"

	"github.com/ve-weiyi/ve-blog-golang/tools/cmd/model/mysql"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "model",
		Short: "生成 Model 模型代码",
	}

	cmd.AddCommand(mysql.NewMysqlCmd())

	return cmd
}
