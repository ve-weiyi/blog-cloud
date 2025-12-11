package mysql

import "github.com/spf13/cobra"

func NewMysqlCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mysql",
		Short: "生成 mysql 数据库模型代码",
	}

	cmd.AddCommand(NewMysqlDDLCmd())
	cmd.AddCommand(NewMysqlDSNCmd())

	return cmd
}
