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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	taskname string
)

//NewLaunchCmd creates Launch command
func NewLaunchCmd(repository *string) *cobra.Command {
	launchCmd := &cobra.Command{
		Use:   "launch",
		Short: "Create a GitHub Source and a Transceiver to automatically generate TaskRuns",
		Run: func(cmd *cobra.Command, args []string) {
			repo = *repository

			if repo != "" {
			    fmt.Printf("%s", GenerateObjBreak(true))
				fmt.Print(GenerateOutput(CreateGithubSource(taskname, repo)))
			    fmt.Printf("%s", GenerateObjBreak(false))
				fmt.Print(GenerateOutput(CreateTransceiver(taskname)))
			    fmt.Printf("%s", GenerateObjLastBreak())
			}
			/*
			TODO handle empty repository way better
			*/
		},
	}
	launchCmd.Flags().StringVarP(&taskname, "task", "t", "", "Task Name to Trigger")
	launchCmd.MarkFlagRequired("task")

	return launchCmd
}

//CreateGithubSource creates Github source based on provided Task name
func CreateGithubSource(taskname string, repo string) sources.GitHubSource {
	var tname = "taskrun-transceiver-"
	tname += taskname
	return sources.GitHubSource{
				TypeMeta: metav1.TypeMeta{
					Kind:       "GitHubSource",
					APIVersion: sources.SchemeGroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: taskname,
				},
				Spec: sources.GitHubSourceSpec{
					OwnerAndRepository : repo,
					EventTypes: []string{"push"},
					AccessToken: sources.SecretValueFromSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "githubsecret",
							},
							Key: "accessToken",
						},
					},
					SecretToken: sources.SecretValueFromSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "githubsecret",
							},
							Key: "secretToken",
						},
					},
					Sink: &corev1.ObjectReference{
							Name:       tname,
							Kind:       "Service",
							APIVersion: "serving.knative.dev/v1alpha1",
					},
				},
	}
}

//CreateTransceiver creates Transceiver object
func CreateTransceiver(taskname string) serving.Service {
	var tname = "taskrun-transceiver-"
	tname += taskname
	return serving.Service{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Service",
				APIVersion: "serving.knative.dev/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: tname,
				Labels: map[string]string{
					"serving.knative.dev/visibility": "cluster-local",
					},
			},
			Spec: serving.ServiceSpec{
				RunLatest: &serving.RunLatestType{
					Configuration: serving.ConfigurationSpec{
						RevisionTemplate: serving.RevisionTemplateSpec{
							Spec: serving.RevisionSpec{
								Container: corev1.Container{
									Image: "gcr.io/triggermesh/transceiver-60a15ebeaf09df9f7ef1bd5f51a22549:latest",
									Env: []corev1.EnvVar{
										{
											Name: "TASK_NAME",
											Value: taskname,
										},
										{
											Name: "NAMESPACE",
											ValueFrom: &corev1.EnvVarSource{
												FieldRef: &corev1.ObjectFieldSelector{
													FieldPath: "metadata.namespace",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
	}
}
