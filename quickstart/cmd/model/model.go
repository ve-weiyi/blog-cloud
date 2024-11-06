/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package model

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "model",
		Short: "从sql文件生成go代码",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
			return
		},
	}

	rootCmd.AddCommand(NewModelDDLCmd().CMD)
	rootCmd.AddCommand(NewModelDSNCmd().CMD)
	return rootCmd
}
