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

	"github.com/spf13/cobra"
	//	corev1 "k8s.io/api/core/v1"
	//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//	sources "github.com/knative/eventing-sources/pkg/apis/sources/v1alpha1"
	//	serving "github.com/knative/serving/pkg/apis/serving/v1alpha1"
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
func CreateGithubSource(taskname string) string {
	return ""
	/*
			return sources.GitHubSource{
				TypeMeta: metav1.TypeMeta{
					Kind:       "GitHubSource",
					APIVersion: sources.SchemeGroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
				Spec: sources.GitHubSourceSpec{
					OwnerAndRepository : "sebgoa/foo",
					EventTypes: []string{"push"},
					AccessToken: sources.SecretValueFromSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "githubsecret",
							},
							Key: "accesstoken",
						},
					},
					SecretToken: sources.SecretValueFromSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "githubsecret",
							},
							Key: "secrettoken",
						},
					},
					Sink: &corev1.ObjectReference{
							Name:       taskname,
							Kind:       "Service",
							APIVersion: "serving.knative.dev/v1alpha1",
				},
			},
		}}
	*/
}

//CreateTransceiver creates Transceiver object
func CreateTransceiver(taskname string) string {
	return ""
	/*
		return serving.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Service",
				APIVersion: "serving.knative.dev/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "foo",
			},
			Spec: serving.ServiceSpec{
		}
	*/
}
