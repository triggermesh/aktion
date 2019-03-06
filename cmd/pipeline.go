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
	"strings"
	"time"

	"github.com/actions/workflow-parser/model"
	pipeline "github.com/knative/build-pipeline/pkg/apis/pipeline/v1alpha1"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	event    string
	repo     string
	registry string
)

type Task struct {
	Identifier string
	Image      string
	Cmd        []string
	Args       []string
	Envs       []corev1.EnvVar
}

type Tasks struct {
	Identifier string
	Task       []Task
}

func NewPipelineCmd() *cobra.Command {
	pipelineCmd := &cobra.Command{
		Use:   "pipeline",
		Short: "Convert the Github Action workflow into a Knative Pipeline task list",
		Run: func(cmd *cobra.Command, args []string) {
			config := ParseData()

			for _, act := range config.Workflows {
				taskRun := CreateTaskRun(act.Identifier)
				tasks := extractTasks(act.Identifier, config)

				GenerateOutput(CreateTask(tasks))
				fmt.Println("---")
				GenerateOutput(taskRun)
			}
		},
	}

	pipelineCmd.Flags().StringVarP(&repo, "repo", "", "", "Upstream git repository")
	pipelineCmd.Flags().StringVarP(&registry, "registry", "r", "knative.registry.svc.cluster.local", "Default docker registry")

	return pipelineCmd
}

func extractTasks(name string, config *model.Configuration) Tasks {
	tasks := Tasks{
		Identifier: name,
		Task:       make([]Task, 0),
	}
	workflow := config.GetWorkflow(name)

	extractedTasks := make([]Task, 0)
	for _, a := range workflow.Resolves {
		extractedTasks = append(extractedTasks, extractActions(config.GetAction(a), config)...)
	}
	tasks.Task = extractedTasks

	return tasks
}

func extractActions(action *model.Action, config *model.Configuration) []Task {
	tasks := make([]Task, 0)

	if len(action.Needs) > 0 {
		for _, a := range action.Needs {
			tasks = append(tasks, extractActions(config.GetAction(a), config)...)
		}
	}

	if action.Uses == nil {
		return tasks
	}

	task := Task{
		Identifier: action.Identifier,
	}

	if !strings.HasPrefix(action.Uses.String(), "docker://") {
		Panic("Can only support docker images for now.")
	}

	task.Image = strings.TrimPrefix(action.Uses.String(), "docker://")

	if action.Runs != nil {
		task.Cmd = action.Runs.Split()
	}

	if action.Args != nil {
		task.Args = action.Args.Split()
	}

	task.Envs = make([]corev1.EnvVar, 0)
	for k, v := range action.Env {
		env := corev1.EnvVar{
			Name:  k,
			Value: v,
		}

		task.Envs = append(task.Envs, env)
	}

	return append(tasks, task)
}

func CreateTaskRun(name string) pipeline.TaskRun {
	taskRun := pipeline.TaskRun{
		Spec: pipeline.TaskRunSpec{
			TaskRef: &pipeline.TaskRef{
				Name: name,
			},
			Trigger: pipeline.TaskTrigger{
				Type: pipeline.TaskTriggerTypeManual,
			},
		},
	}

	taskRun.SetDefaults()
	taskRun.TypeMeta = metav1.TypeMeta{
		Kind: "TaskRun",
		APIVersion: "pipeline.knative.dev/v1alpha1",
	}

	taskRun.ObjectMeta = metav1.ObjectMeta{
		Name: name,
		CreationTimestamp: metav1.Time{time.Now()},
	}

	err := taskRun.Validate()
	if err != nil {
		Panic("Failed validation: %s\n", err)
	}

	return taskRun
}

func CreateTask(tasks Tasks) pipeline.Task {
	task := pipeline.Task{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Task",
			APIVersion: "pipeline.knative.dev/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: tasks.Identifier,
		},
	}

	var taskSpec pipeline.TaskSpec
	steps := make([]corev1.Container, 0)
	for _, t := range tasks.Task {
		steps = append(steps, createContainer(t))
	}
	taskSpec.Steps = steps
	task.Spec = taskSpec

	return task
}

func createContainer(task Task) corev1.Container {
	return corev1.Container{
		Name:    convertName(task.Identifier),
		Image:   task.Image,
		Command: task.Cmd,
		Args:    task.Args,
		Env:     task.Envs,
	}
}

func convertName(name string) string {
	n := strings.Replace(name, " ", "-", -1)
	return strings.ToLower(n)
}
