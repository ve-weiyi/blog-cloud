/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "blog-gin",
		Short:   "ve-blog-golang 博客系统",
		Version: fmt.Sprintf("%s %s/%s", "v1.0.0", runtime.GOOS, runtime.GOARCH),
		Run: func(cmd *cobra.Command, args []string) {
			//实现功能逻辑的函数。
			//_ = cmd.Help()
			NewApiCmd().Execute()
		},
	}

	rootCmd.AddCommand(NewApiCmd())
	rootCmd.AddCommand(NewMigrateCmd())

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
