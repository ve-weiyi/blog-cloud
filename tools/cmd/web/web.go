package web

import (
	"github.com/spf13/cobra"

	"github.com/ve-weiyi/ve-blog-golang/tools/cmd/web/ts"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "web",
		Short: "生成 Web 前端代码",
	}

	cmd.AddCommand(ts.NewTypescriptCmd())
	return cmd
}
