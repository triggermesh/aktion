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
	"github.com/spf13/cobra"

	"github.com/triggermesh/aktion/pkg/client"

	pipeline "github.com/knative/build-pipeline/pkg/apis/pipeline/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	registry                string
	revision 				string
	taskrun                 bool
	visitedActionDependency map[string]bool
	pipelineResources       map[string]*Image
	applyPipelineFlag       bool
)

type ImageConst int

const (
	DOCKER ImageConst = iota
	GIT
	LOCAL
)

type Image struct {
	Type ImageConst
	Path string
	BuildTaskName string
	PipelineResourceSource pipeline.PipelineResource
	PipelineResourceImage pipeline.PipelineResource
}

//Task represents Task object
type Task struct {
	Identifier string
	Image      *Image
	Cmd        []string
	Args       []string
	Envs       []corev1.EnvVar
	EnvFrom    []corev1.EnvFromSource
}

//Tasks groups Task objects by one Identifier
type Tasks struct {
	Identifier string
	Task       []Task
}

//NewCreateCmd creates new create command
func NewCreateCmd(kubeConfig *string, ns *string, repository *string) *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Convert the Github Action workflow into a Tekton Task list",
		Run: func(cmd *cobra.Command, args []string) {
			config := ParseData()
			visitedActionDependency = make(map[string]bool)
			pipelineResources = make(map[string]*Image)
			namespace = *ns
			repo = *repository

			/* CAB: This will need to be refactored to account for pipeline revamping */
			/*
			if repo != "" {
				repoPipelineResource := createPipelineResource(repo, config)

				fmt.Println("---")
				fmt.Print(GenerateOutput(repoPipelineResource))
				fmt.Println("---")
			}
			*/

			for _, act := range config.Workflows {
				tasks := extractTasks(act.Identifier, config)
				primaryPipeline := CreatePipeline(tasks, act.Identifier, repo)
				pipelineRun := CreatePipelineRun(act.Identifier)

				if applyPipelineFlag {
					applyPipeline(*kubeConfig, primaryPipeline, pipelineRun, CreateTask(tasks, repo))
				} else {
					for _,v := range pipelineResources {
						fmt.Printf("%s", GenerateOutput(createPipelineResource(*v, true)))
						fmt.Printf("---\n")
						fmt.Printf("%s", GenerateOutput(createPipelineResource(*v, false)))
						fmt.Printf("---\n")
						fmt.Printf("%s", GenerateOutput(CreateBuildTask(*v)))
						fmt.Printf("---\n")
					}

					fmt.Printf("%s", GenerateOutput(CreateTask(tasks, repo)))

					fmt.Printf("---\n")
					fmt.Printf("%s", GenerateOutput(primaryPipeline))

					if taskrun {
						//fmt.Printf("---\n%s", GenerateOutput(taskRun))
						fmt.Printf("---\n%s", GenerateOutput(pipelineRun))
					}
				}
			}
		},
	}

	createCmd.Flags().StringVarP(&repo, "repo", "", "", "Upstream git repository")
	createCmd.Flags().StringVarP(&revision, "revision", "", "master", "Upstream repository revision, branch, or tag")
	createCmd.Flags().StringVarP(&registry, "registry", "r", "http://knative.registry.svc.cluster.local", "Default docker registry")
	createCmd.Flags().BoolVarP(&taskrun, "taskrun", "t", false, "Flag to create TaskRun")
	createCmd.Flags().BoolVarP(&applyPipelineFlag, "apply", "a", false, "Apply the generated Tekton pipeline to the user's kubernetes cluster")

	return createCmd
}

// will need to add a lot more to the generation
func applyPipeline(kubeConfig string, primaryPipeline pipeline.Pipeline, pipelineRun pipeline.PipelineRun, tasks pipeline.Task) {
	// add if check for taskrun to build/inject the task
	clientSet, err := client.NewClient(client.ConfigPath(kubeConfig))
	if err != nil {
		Panic("Error connecting to kubernetes cluster: %s\n", err)
	}

	// Apply resources
	for _, v := range pipelineResources {
		srcResource := createPipelineResource(*v, true)
		imgResource := createPipelineResource(*v, false)
		buildTask := CreateBuildTask(*v)

		_, err = clientSet.Pipeline.TektonV1alpha1().PipelineResources(namespace).Create(&srcResource)
		if err != nil {
			Panic("Unable to create source reference resource: %s\n", err)
		}

		_, err = clientSet.Pipeline.TektonV1alpha1().PipelineResources(namespace).Create(&imgResource)
		if err != nil {
			Panic("Unable to create source reference image: %s\n", err)
		}

		_, err = clientSet.Pipeline.TektonV1alpha1().Tasks(namespace).Create(&buildTask)
		if err != nil {
			Panic("Unable to create build task: %s\n", err)
		}
	}

	_, err = clientSet.Pipeline.TektonV1alpha1().Tasks(namespace).Create(&tasks)
	if err != nil {
		Panic("Unable to create tasks: %s\n", err)
	}

	_, err = clientSet.Pipeline.TektonV1alpha1().Pipelines(namespace).Create(&primaryPipeline)
	if err != nil {
		Panic("Unable to create primary pipeline: %s\n", err)
	}

	if taskrun {
		_, err = clientSet.Pipeline.TektonV1alpha1().PipelineRuns(namespace).Create(&pipelineRun)
		if err != nil {
			Panic("Unable to create pipeline run: %s\n", err)
		}
	}
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
			if !visitedActionDependency[config.GetAction(a).Identifier] {
				tasks = append(tasks, extractActions(config.GetAction(a), config)...)
			}
		}
	}

	if action.Uses == nil {
		return tasks
	}

	task := Task{
		Identifier: action.Identifier,
	}

	if strings.HasPrefix(action.Uses.String(), "docker://") {
		task.Image = &Image{
			Type: DOCKER,
			Path: strings.TrimPrefix(action.Uses.String(), "docker://"),
		}
	} else if strings.HasPrefix(action.Uses.String(), "./") {
		if len(repo) == 0 {
			Panic("The repo flag must be specified to use the action: %s\n", action.Identifier)
		} else if pipelineResources[convertUsesName(action.Uses.String())] != nil {
			task.Image = pipelineResources[convertUsesName(action.Uses.String())]
		} else {
			task.Image = &Image{
				Type: LOCAL,
				Path: strings.TrimPrefix(action.Uses.String(), "./"),
				BuildTaskName: convertUsesName(action.Uses.String()),
			}

			task.Image.PipelineResourceSource = createPipelineResource(*task.Image, true)
			task.Image.PipelineResourceImage = createPipelineResource(*task.Image, false)
			pipelineResources[convertUsesName(action.Uses.String())] = task.Image
		}
	} else if strings.Contains(action.Uses.String(), "@") {
		/* TODO: Should this use the repo flag, or do we need a new flag to reflect another repo? */
		if pipelineResources[convertUsesName(action.Uses.String())] != nil {
			task.Image = pipelineResources[convertUsesName(action.Uses.String())]
		} else {
			task.Image = &Image{
				Type: GIT,
				Path: "github.com/" + action.Uses.String(),
				BuildTaskName: convertUsesName(action.Uses.String()),
			}

			task.Image.PipelineResourceSource = createPipelineResource(*task.Image, true)
			task.Image.PipelineResourceImage = createPipelineResource(*task.Image, false)
			pipelineResources[convertUsesName(action.Uses.String())] = task.Image
		}
	} else {
		Panic("The image %s for %s is unsupported\n", action.Uses.String(), action.Identifier)
	}

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

	if action.Secrets != nil {
		task.EnvFrom = make([]corev1.EnvFromSource, 0)
		for _, s := range action.Secrets {
			secret := corev1.EnvFromSource{
				SecretRef: &corev1.SecretEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: s}},
			}
			task.EnvFrom = append(task.EnvFrom, secret)
		}
	}

	// Mark as visited
	visitedActionDependency[task.Identifier] = true

	return append(tasks, task)
}

// CreateTaskRun function creates TaskRun object - This is considered a low-level
// operation. Use PipelineRun instead
func CreateTaskRun(name string) pipeline.TaskRun {
	taskRun := pipeline.TaskRun{
		Spec: pipeline.TaskRunSpec{
			TaskRef: &pipeline.TaskRef{
				Name: convertName(name),
			},
			Trigger: pipeline.TaskTrigger{
				Type: pipeline.TaskTriggerTypeManual,
			},
		},
	}

	taskRun.SetDefaults()
	taskRun.TypeMeta = metav1.TypeMeta{
		Kind:       "TaskRun",
		APIVersion: "tekton.dev/v1alpha1",
	}

	taskRun.ObjectMeta = metav1.ObjectMeta{
		Name:              convertName(name),
		CreationTimestamp: metav1.Time{time.Now()},
	}

	err := taskRun.Validate()
	if err != nil {
		Panic("Failed validation: %s\n", err)
	}

	return taskRun
}

// CreatePipeline Generates the pipeline and associated tasks
func CreatePipeline(tasks Tasks, name string, repo string) pipeline.Pipeline {
	line := pipeline.Pipeline{
		TypeMeta: metav1.TypeMeta{
			Kind:	"Pipeline",
			APIVersion: "tekton.dev/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              convertName(name + "-pipeline"),
			CreationTimestamp: metav1.Time{time.Now()},
		},
	}

	// TODO: Do we want to apply the custom resource definition for defined repos?

	specResources := make([]pipeline.PipelineDeclaredResource, 0)
	specTasks := make([]pipeline.Task, 0)
	specPipelineTask := make([]pipeline.PipelineTask, 0)

	for _, v  := range pipelineResources {
		srcResource := pipeline.PipelineDeclaredResource{
			Name: v.PipelineResourceSource.ObjectMeta.Name,
			Type: v.PipelineResourceSource.Spec.Type,
		}

		imgResource := pipeline.PipelineDeclaredResource{
			Name: v.PipelineResourceImage.ObjectMeta.Name,
			Type: v.PipelineResourceImage.Spec.Type,
		}

		specResources = append(specResources, srcResource)
		specResources = append(specResources, imgResource)
		buildTask := CreateBuildTask(*v)
		specTasks = append(specTasks, buildTask)

		pipelineBuildTask := pipeline.PipelineTask{
			Name: "pipeline-build-" + buildTask.Name,
			TaskRef: pipeline.TaskRef{
				Name: buildTask.Name,
			},
			Resources: &pipeline.PipelineTaskResources{
				Inputs: []pipeline.PipelineTaskInputResource{{
					Name: "workspace",
					Resource: v.PipelineResourceSource.ObjectMeta.Name,
				}},
				Outputs: []pipeline.PipelineTaskOutputResource{{
					Name: "image",
					Resource: v.PipelineResourceImage.ObjectMeta.Name,
				}},
			},
			Params: []pipeline.Param{{
				Name: "pathToContext",
				Value: "/workspace/" + convertName(name) + "/" + extractRepoPath(v.Path),
			}},
		}

		specPipelineTask = append(specPipelineTask, pipelineBuildTask)
	}

	task := CreateTask(tasks, repo)
	primaryPipelineTask := pipeline.PipelineTask{
		Name: task.Name,
		TaskRef: pipeline.TaskRef{
			Name: task.Name,
		},
	}
	specPipelineTask = append(specPipelineTask, primaryPipelineTask)

	line.Spec.Resources = specResources
	line.Spec.Tasks = specPipelineTask

	return line
}

func CreatePipelineRun(name string) pipeline.PipelineRun {
	pipelineRun := pipeline.PipelineRun{
		Spec: pipeline.PipelineRunSpec{
			PipelineRef: pipeline.PipelineRef{
				Name: convertName(name + "-pipeline"),
			},
			Trigger: pipeline.PipelineTrigger{
				Type: pipeline.PipelineTriggerTypeManual,
			},
		},
	}

	pipelineRun.TypeMeta = metav1.TypeMeta{
		Kind:	"PipelineRun",
		APIVersion: "tekton.dev/v1alpha1",
	}

	pipelineRun.ObjectMeta = metav1.ObjectMeta{
		Name:              convertName(name + "-pipeline-run"),
		CreationTimestamp: metav1.Time{time.Now()},
	}

	err := pipelineRun.Validate()
	if err != nil {
		Panic("Failed validation for pipeline-run: %s\n", err)
	}

	return pipelineRun
}

//CreateTask creates Task object
func CreateTask(tasks Tasks, repo string) pipeline.Task {
	task := pipeline.Task{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Task",
			APIVersion: "tekton.dev/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: convertName(tasks.Identifier),
		},
	}

	var taskSpec pipeline.TaskSpec
	steps := make([]corev1.Container, 0)

	if repo != "" {
		taskSpec.Inputs = &pipeline.Inputs{
			Resources: []pipeline.TaskResource{{
				Name: convertName(tasks.Identifier),
				Type: "git",
			}},
		}
	}

	for _, t := range tasks.Task {
		steps = append(steps, createContainer(t))
	}
	taskSpec.Steps = steps
	task.Spec = taskSpec

	return task
}

// Given the github-action repo designation of org/repo/path..., return just the org/repo portion
func extractRepoPrefix(repo string) string {
	path := strings.Split(repo, "@")[0]

	if strings.Count(path, "/") == 1 {
		return path
	}

	components := strings.Split(repo, "/")

	return components[0] + "/" + strings.Split(components[1], "@")[0]
}

// Given the github-action repo designation of org/repo/path..., return just the path portion
func extractRepoPath(repo string) string {
	path := strings.Split(repo, "@")[0]

	if strings.Count(path, "/") == 1 {
		return "/"
	}

	return strings.TrimPrefix(path, extractRepoPrefix(path))
}

// Extract the provided hash, branch, tag. Default to "master" if one is not provided
func extractRepoRevision(repo string) string {
	rev := strings.Split(repo, "@")

	if rev[0] == repo {
		return "master"
	}

	return rev[1]
}

/*
 * createPipelineResource will generate a new resource object based off the Image
 * contents.
 *
 * resourceType - true indicates an Input resource (git), false indicates an Output resource (image)
 */
func createPipelineResource(image Image, resourceType bool) pipeline.PipelineResource {
	resourceName := image.BuildTaskName

	if resourceType {
		resourceName += "-git"
	} else {
		resourceName += "-image"
	}

	resource := pipeline.PipelineResource{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PipelineResource",
			APIVersion: "tekton.dev/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: convertName(resourceName),
		},
	}

	resourceParams := make([]pipeline.Param, 0)

	if resourceType {
		// There are two additional fields that need to be extracted from the
		// workflow representation of the repo: The URL and the revision.
		var url string

		if image.Type == LOCAL {
			hasVersion := strings.IndexAny(repo, "@")
			if hasVersion == -1 {
				url = repo
				revision = "master"
			} else {
				components := strings.Split(repo, "@")
				url = components[0]
				revision = components[1]
			}
		} else if image.Type == GIT {
			// TODO: If repo is passed as an argument, do we use that to override this?
			url = "https://github.com/" + extractRepoPrefix(image.Path)
			revision = extractRepoRevision(image.Path) // This is for the 3rd party repo being accessed
		}

		resourceParams = append(resourceParams,
			pipeline.Param{
				Name: "revision",
				Value: revision,
		})

		resourceParams = append(resourceParams,
			pipeline.Param{
				Name: "url",
				Value: url,
		})

		resource.Spec = pipeline.PipelineResourceSpec{
			Type: pipeline.PipelineResourceTypeGit,
			Params: resourceParams,
		}
	} else {
		resourceParams = append(resourceParams,
			pipeline.Param{
				Name: "url",
				Value: registry + "/" + image.BuildTaskName,
		})

		resource.Spec = pipeline.PipelineResourceSpec{
			Type: pipeline.PipelineResourceTypeImage,
			Params: resourceParams,
		}
	}

	return resource
}

// CreateBuildTask create a task to clone a git repo and use Kaniko to build the docker image
func CreateBuildTask(image Image) pipeline.Task {
	task := pipeline.Task{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Task",
			APIVersion: "tekton.dev/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "build-" + convertName(image.BuildTaskName),
		},
	}

	// CAB: The pathToDocker and pathToContext get set in the pipeline when calling taskRef
	task.Spec = pipeline.TaskSpec{
		Inputs: &pipeline.Inputs{
			Resources: []pipeline.TaskResource{{
				Name: image.BuildTaskName + "-git",
				Type: pipeline.PipelineResourceTypeGit,
			}},
			Params: []pipeline.TaskParam{
				{
					Name: "pathToDockerFile",
					Default: "Dockerfile",
				},
				{
					Name: "pathToContext",
				},
			},
		},
		Outputs: &pipeline.Outputs{
			Resources: []pipeline.TaskResource{{
				Name: image.BuildTaskName + "-image",
				Type: pipeline.PipelineResourceTypeImage,
			}},
		},
		Steps: []corev1.Container{{
			Name:    "build-and-push-" + convertName(image.BuildTaskName),
			Image:   "gcr.io/kaniko-project/executor",
			Command: []string{"/kaniko/executor"},
			Args: []string{
				"--dockerfile=${inputs.params.pathToDockerFile}",
				"--destination=${outputs.resources.builtImage.url}",
				"--context=${inputs.params.pathToContext}",
			},
		}},
	}

	return task
}

func createContainer(task Task) corev1.Container {
	return corev1.Container{
		Name:    convertName(task.Identifier),
		Image:   task.Image.Path,
		Command: task.Cmd,
		Args:    task.Args,
		Env:     task.Envs,
		EnvFrom: task.EnvFrom,
	}
}

func convertName(name string) string {
	n := strings.Replace(name, " ", "-", -1)
	return strings.ToLower(n)
}

// Convert the workflow Uses entry to a common name that can be referenced
func convertUsesName(name string) string {
	n := strings.Split(name, "@")[0]
	n = strings.Replace(n, "/", "-", -1)
	n = strings.Replace(n, ".", "-", -1)
	n = strings.Replace(n, "--", "-", -1)
	n = strings.TrimPrefix(n, "-")

	return strings.ToLower(n)
}
