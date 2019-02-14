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
    "gopkg.in/yaml.v2"
)

func NewParserCmd(filename *string, outputType *string) *cobra.Command {
	parserCmd := &cobra.Command{
		Use:   "parser",
		Short: "Parse the workflow into a JSON file",
		Run: func(cmd *cobra.Command, args []string) {
			f, err := os.Open(*filename)
			var output []byte

			if err != nil {
				Panic("Error opening file: %s\n", err)
			}

			config, err := parser.Parse(f)
			if err != nil {
				Panic("Error parsing file: %s\n", err)
			}
			f.Close()

			if (*outputType == "json") {
				output, err = json.MarshalIndent(config, "", "  ")
				if err != nil {
					Panic("Error generating JSON output: %s\n", err)
				}
			} else if (*outputType == "yaml") {
				output, err = yaml.Marshal(config)
				if err != nil {
					Panic("Error generating YAML output: %s\n", err)
				}
			} else {
				Panic("Unsupported format: %s. Expect json or yaml\n", *outputType)
			}

			fmt.Printf("%s\n", output)
		},
	}

	return parserCmd
}
