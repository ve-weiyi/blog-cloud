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
	"os"

	"github.com/spf13/cobra"
)

type RootCmd struct {
	cmd *cobra.Command
}

func NewRootCmd() *RootCmd {
	var rootCmd = &cobra.Command{
		Use:   "",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			//实现功能逻辑的函数。
			_ = cmd.Help()
			return
		},
	}

	root := &RootCmd{
		cmd: rootCmd,
	}

	root.cmd.AddCommand(NewServerCmd().cmd)
	root.cmd.AddCommand(NewMigrateCmd().cmd)
	root.cmd.AddCommand(NewRabbitmqCmd().cmd)
	root.init()
	return root
}

func (s *RootCmd) init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.server.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

}

func (s *RootCmd) Execute() {
	err := s.cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
