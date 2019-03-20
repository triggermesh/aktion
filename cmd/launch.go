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

	sources "github.com/knative/eventing-sources/pkg/apis/sources/v1alpha1"
	serving "github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"github.com/spf13/cobra"
)

var (
	taskname string
)

//NewLaunchCmd creates Launch command
func NewLaunchCmd() *cobra.Command {
	launchCmd := &cobra.Command{
		Use:   "launch",
		Short: "launch whatever",
		Run: func(cmd *cobra.Command, args []string) {
			GenerateOutput(CreateGithubSource(taskname))
			fmt.Println("---")
			GenerateOutput(CreateTransceiver(taskname))
			fmt.Println("---")
		},
	}
	launchCmd.Flags().StringVarP(&taskname, "taskname", "t", "", "Task Name to Trigger")

	return launchCmd
}

//CreateGithubSource creates Github source based on provided Task name
func CreateGithubSource(taskname string) sources.GitHubSource {
	return sources.GitHubSource{}
}

//CreateTransceiver creates Transceiver object
func CreateTransceiver(taskname string) serving.Service {
	return serving.Service{}
}
