package ts

import (
	"github.com/spf13/cobra"
)

func NewTypescriptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ts",
		Short: "生成 TypeScript 代码",
	}

	cmd.AddCommand(NewTypescriptApiCmd())
	cmd.AddCommand(NewTypescriptSwaggerCmd())
	return cmd
}
