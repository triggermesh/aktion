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
	"encoding/json"
	"fmt"
	"os"

	"github.com/actions/workflow-parser/model"
	"github.com/actions/workflow-parser/parser"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

var (
	version    string
	filename   string
	outputType string
	kubeConfig string
	namespace  string
	repo       string
)

var aktionCmd = &cobra.Command{
	Use:     "aktion",
	Short:   "Convert GitHub Actions workflow into Tekton resources",
	Version: version,
}

//Panic is a wrapper to format pannic messages
func Panic(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

//GenerateOutput turns any data of type json or yaml into a string, panic otherwise
func GenerateOutput(data interface{}) string {
	var output []byte
	var err error

	if outputType == "json" {
		output, err = json.MarshalIndent(data, "", "  ")
		if err != nil {
			Panic("Error generating JSON output: %s\n", err)
		}
	} else if outputType == "yaml" {
		output, err = yaml.Marshal(data)
		if err != nil {
			Panic("Error generating YAML output: %s\n", err)
		}
	} else {
		Panic("Unsupported format: %s. Expect json or yaml\n", outputType)
	}

	return fmt.Sprintf("%s", output)
}

//ParseData parses Github Action Workflow File into Configuration object
func ParseData() *model.Configuration {
	f, err := os.Open(filename)

	if err != nil {
		Panic("Error opening file: %s\n", err)
	}

	config, err := parser.Parse(f)
	if err != nil {
		Panic("Error parsing file: %s\n", err)
	}
	f.Close()

	return config
}

//Execute launches aktion command
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
	aktionCmd.PersistentFlags().StringVarP(&outputType, "output", "o", "yaml", "Output type for the results (json|yaml)")
	aktionCmd.PersistentFlags().StringVarP(&kubeConfig, "kubeconfig", "k", "", "Kubernetes config file")
	aktionCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes namespace")
	aktionCmd.PersistentFlags().StringVarP(&repo, "git", "g", "", "Git repository")
	aktionCmd.AddCommand(versionCmd)
	aktionCmd.AddCommand(NewParserCmd())
	aktionCmd.AddCommand(NewCreateCmd(&kubeConfig, &namespace, &repo))
	aktionCmd.AddCommand(NewLaunchCmd(&repo))
}

func initConfig() {
	// do nothing for now
}
