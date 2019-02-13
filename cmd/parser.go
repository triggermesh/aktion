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

	"github.com/actions/workflow-parser/parser"
	"github.com/spf13/cobra"
)

func NewParserCmd(filename *string) *cobra.Command {
	return &cobra.Command{
		Use:   "parser",
		Short: "Parse the workflow into a JSON file",
		Run: func(cmd *cobra.Command, args []string) {
			f, err := os.Open(*filename)
			if err != nil {
				fmt.Printf("Error opening file (%s): %s\n", *filename, err)
				os.Exit(1)
			}

			config, err := parser.Parse(f)
			if err != nil {
				fmt.Printf("Error parsing file (%s): %s\n", *filename, err)
				os.Exit(1)
			}

			jsonOutput, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				fmt.Printf("Error generating JSON output: %s\n", err)
				os.Exit(1)
			}

			fmt.Printf("%s\n", jsonOutput)

			f.Close()
		},
	}
}
