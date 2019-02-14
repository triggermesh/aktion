/*
Copyright (c) 2019 TriggerMesh, Inc

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

	"github.com/spf13/cobra"
)

var (
	version  string
	filename string
	outputType string
)

var aktionCmd = &cobra.Command{
	Use:     "aktion",
	Short:   "Actions for Knative",
	Version: version,
}

func Panic(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func Execute() {
	if err := aktionCmd.Execute(); err != nil {
		Panic("Error: %s\n", err)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s, version %s\n", aktionCmd.Short, version)
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	aktionCmd.PersistentFlags().StringVarP(&filename, "filename", "f", "main.workflow", "Github Action Workflow File")
	aktionCmd.PersistentFlags().StringVarP(&outputType, "outputType", "o", "json", "Output type for the results (json|yaml)")
	aktionCmd.AddCommand(versionCmd)
	aktionCmd.AddCommand(NewParserCmd(&filename, &outputType))
}

func initConfig() {
	// do nothing for now
}
