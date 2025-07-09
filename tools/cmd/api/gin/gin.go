package gin

import "github.com/spf13/cobra"

func NewGinCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gin",
		Short: "生成 Gin 框架代码",
	}

	cmd.AddCommand(NewGinApiCmd())
	cmd.AddCommand(NewGinSwaggerCmd())

	return cmd
}
