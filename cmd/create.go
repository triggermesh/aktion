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
	revision                string
	pipelinerun             bool
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
	Type                   ImageConst
	Path                   string
	BuildTaskName          string
	PipelineResourceSource pipeline.PipelineResource
	PipelineResourceImage  pipeline.PipelineResource
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
func NewCreateCmd(kubeConfig *string, ns *string, gitRepository *string) *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Convert the Github Action workflow into a Tekton Task list",
		Run: func(cmd *cobra.Command, args []string) {
			config := ParseData()
			visitedActionDependency = make(map[string]bool)
			pipelineResources = make(map[string]*Image)
			namespace = *ns
			repo = *gitRepository

			for _, act := range config.Workflows {
				tasks := extractTasks(act.Identifier, config)
				primaryPipeline := createPipeline(tasks, act.Identifier, repo)
				pipelineRun := createPipelineRun(act.Identifier, repo, act.Identifier)
				pipelineRepo := createRepoPipelineResource(repo, act.Identifier, config)

				if applyPipelineFlag {
					applyPipeline(*kubeConfig, primaryPipeline, pipelineRun, createTask(tasks, repo), pipelineRepo)
				} else {
					fmt.Printf("%s", GenerateObjBreak(true))

					for _, v := range pipelineResources {
						fmt.Printf("%s", GenerateOutput(createPipelineResource(*v, true)))
						fmt.Printf("%s", GenerateObjBreak(false))
						fmt.Printf("%s", GenerateOutput(createPipelineResource(*v, false)))
						fmt.Printf("%s", GenerateObjBreak(false))
						fmt.Printf("%s", GenerateOutput(createBuildTask(*v)))
						fmt.Printf("%s", GenerateObjBreak(false))
					}

					if repo != "" {
						fmt.Printf("%s", GenerateOutput(*pipelineRepo))
						fmt.Printf("%s", GenerateObjBreak(false))
					}

					fmt.Printf("%s", GenerateOutput(createTask(tasks, repo)))

					fmt.Printf("%s", GenerateObjBreak(false))
					fmt.Printf("%s", GenerateOutput(primaryPipeline))

					if pipelinerun {
						//fmt.Printf("---\n%s", GenerateOutput(taskRun))
						fmt.Printf("%s", GenerateObjBreak(false))
						fmt.Printf("%s", GenerateOutput(pipelineRun))
					}
				}
			}

			fmt.Printf("%s", GenerateObjLastBreak())
		},
	}

	createCmd.Flags().StringVarP(&revision, "revision", "", "master", "Upstream repository revision, branch, or tag")
	createCmd.Flags().StringVarP(&registry, "registry", "r", "knative.registry.svc.cluster.local", "Default docker registry")
	createCmd.Flags().BoolVarP(&pipelinerun, "pipelinerun", "p", false, "Flag to create PipelineRun")
	createCmd.Flags().BoolVarP(&applyPipelineFlag, "apply", "a", false, "Apply the generated Tekton pipeline to the user's kubernetes cluster")

	return createCmd
}

// will need to add a lot more to the generation
func applyPipeline(kubeConfig string, primaryPipeline pipeline.Pipeline, pipelineRun pipeline.PipelineRun, tasks pipeline.Task, pipelineRepo *pipeline.PipelineResource) {
	// add if check for pipelinerun to build/inject the task
	clientSet, err := client.NewClient(client.ConfigPath(kubeConfig))
	if err != nil {
		Panic("Error connecting to kubernetes cluster: %s\n", err)
	}

	// Apply resources
	for _, v := range pipelineResources {
		srcResource := createPipelineResource(*v, true)
		imgResource := createPipelineResource(*v, false)
		buildTask := createBuildTask(*v)

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

	// Apply global git repo
	if pipelineRepo != nil {
		_, err = clientSet.Pipeline.TektonV1alpha1().PipelineResources(namespace).Create(pipelineRepo)
		if err != nil {
			Panic("Unable to create global repo resource: %s\n", err)
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

	if pipelinerun {
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
			Panic("The git flag must be specified to use the action: %s\n", action.Identifier)
		} else if pipelineResources[convertUsesName(action.Uses.String())] != nil {
			task.Image = pipelineResources[convertUsesName(action.Uses.String())]
		} else {
			task.Image = &Image{
				Type:          LOCAL,
				Path:          strings.TrimPrefix(action.Uses.String(), "./"),
				BuildTaskName: convertUsesName(action.Uses.String()),
			}

			task.Image.PipelineResourceSource = createPipelineResource(*task.Image, true)
			task.Image.PipelineResourceImage = createPipelineResource(*task.Image, false)
			pipelineResources[convertUsesName(action.Uses.String())] = task.Image
		}
	} else if strings.Contains(action.Uses.String(), "@") {
		if pipelineResources[convertUsesName(action.Uses.String())] != nil {
			task.Image = pipelineResources[convertUsesName(action.Uses.String())]
		} else {
			task.Image = &Image{
				Type:          GIT,
				Path:          "github.com/" + action.Uses.String(),
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

// createTaskRun function creates TaskRun object - This is considered a low-level
// operation. Use PipelineRun instead
func createTaskRun(name string) pipeline.TaskRun {
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

// createPipeline Generates the pipeline and associated tasks
func createPipeline(tasks Tasks, name string, repo string) pipeline.Pipeline {
	line := pipeline.Pipeline{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pipeline",
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

	for _, v := range pipelineResources {
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
		buildTask := createBuildTask(*v)
		specTasks = append(specTasks, buildTask)

		pipelineBuildTask := pipeline.PipelineTask{
			Name: buildTask.Name,
			TaskRef: pipeline.TaskRef{
				Name: buildTask.Name,
			},
			Resources: &pipeline.PipelineTaskResources{
				Inputs: []pipeline.PipelineTaskInputResource{{
					Name:     "workspace",
					Resource: v.PipelineResourceSource.ObjectMeta.Name,
				}},
				Outputs: []pipeline.PipelineTaskOutputResource{{
					Name:     "image",
					Resource: v.PipelineResourceImage.ObjectMeta.Name,
				}},
			},
			Params: []pipeline.Param{{
				Name:  "pathToContext",
				Value: "/workspace/" + convertName(name) + "/" + extractRepoPath(v.Path),
			}},
		}

		specPipelineTask = append(specPipelineTask, pipelineBuildTask)
	}

	task := createTask(tasks, repo)
	primaryPipelineTask := pipeline.PipelineTask{
		Name: task.Name,
		TaskRef: pipeline.TaskRef{
			Name: task.Name,
		},
	}

	if repo != "" {
		specResources = append(specResources, pipeline.PipelineDeclaredResource{
			Name: convertName(name),
			Type: pipeline.PipelineResourceTypeGit,
		})

		primaryPipelineTask.Resources = &pipeline.PipelineTaskResources{
			Inputs: []pipeline.PipelineTaskInputResource{{
				Name:     convertName(name),
				Resource: convertName(name),
			}},
		}
	}

	specPipelineTask = append(specPipelineTask, primaryPipelineTask)
	line.Spec.Resources = specResources
	line.Spec.Tasks = specPipelineTask

	return line
}

func createPipelineRun(name string, repo string, workflowName string) pipeline.PipelineRun {
	// setup the resource run bindings
	resourceBindings := make([]pipeline.PipelineResourceBinding, 0)

	for _, v := range pipelineResources {
		resourceBindings = append(resourceBindings, pipeline.PipelineResourceBinding{
			Name: v.PipelineResourceImage.Name,
			ResourceRef: pipeline.PipelineResourceRef{
				Name: v.PipelineResourceImage.Name,
			},
		})

		resourceBindings = append(resourceBindings, pipeline.PipelineResourceBinding{
			Name: v.PipelineResourceSource.Name,
			ResourceRef: pipeline.PipelineResourceRef{
				Name: v.PipelineResourceSource.Name,
			},
		})
	}

	if repo != "" {
		resourceBindings = append(resourceBindings, pipeline.PipelineResourceBinding{
			Name: convertName(workflowName),
			ResourceRef: pipeline.PipelineResourceRef{
				Name: convertName(workflowName),
			},
		})
	}

	pipelineRun := pipeline.PipelineRun{
		Spec: pipeline.PipelineRunSpec{
			PipelineRef: pipeline.PipelineRef{
				Name: convertName(name + "-pipeline"),
			},
			Trigger: pipeline.PipelineTrigger{
				Type: pipeline.PipelineTriggerTypeManual,
			},
			Resources: resourceBindings,
		},
	}

	pipelineRun.TypeMeta = metav1.TypeMeta{
		Kind:       "PipelineRun",
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

//createTask creates Task object
func createTask(tasks Tasks, repo string) pipeline.Task {
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
	basedir := strings.Split(repo, "@")[0]

	/*
		if strings.Count(path, "/") == 1 {
			return path
		}

		components := strings.Split(path, "/")

		// FIXME: Array index out of range

		return components[0] + "/" + strings.Split(components[1], "@")[0]
	*/

	return basedir
}

// Given the github-action repo designation of org/repo/path..., return just the path portion
func extractRepoPath(repo string) string {
	path := strings.Split(repo, "@")[0]

	if strings.Count(path, "/") == 1 {
		return ""
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
			url = "https://" + extractRepoPrefix(image.Path)
			revision = extractRepoRevision(image.Path) // This is for the 3rd party repo being accessed
		}

		resourceParams = append(resourceParams,
			pipeline.Param{
				Name:  "revision",
				Value: revision,
			})

		resourceParams = append(resourceParams,
			pipeline.Param{
				Name:  "url",
				Value: url,
			})

		resource.Spec = pipeline.PipelineResourceSpec{
			Type:   pipeline.PipelineResourceTypeGit,
			Params: resourceParams,
		}
	} else {
		resourceParams = append(resourceParams,
			pipeline.Param{
				Name:  "url",
				Value: registry + "/" + image.BuildTaskName,
			})

		resource.Spec = pipeline.PipelineResourceSpec{
			Type:   pipeline.PipelineResourceTypeImage,
			Params: resourceParams,
		}
	}

	return resource
}

// createRepoPipelineResource - Ensure that we create a PipelineResource object for the CLI specified repo
func createRepoPipelineResource(repo string, workflow string, config *model.Configuration) *pipeline.PipelineResource {
	var url string
	var revision string

	if repo == "" {
		return nil
	}

	// Hack: Get the first worklow in the list to get a name
	resource := pipeline.PipelineResource{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PipelineResource",
			APIVersion: "tekton.dev/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: convertName(workflow),
		},
	}

	hasVersion := strings.IndexAny(repo, "@")
	if hasVersion == -1 {
		url = repo
		revision = "master"
	} else {
		components := strings.Split(repo, "@")
		url = components[0]
		revision = components[1]
	}

	inputparams := make([]pipeline.Param, 0)

	inputparams = append(inputparams, pipeline.Param{
		Name:  "revision",
		Value: revision,
	})

	inputparams = append(inputparams, pipeline.Param{
		Name:  "url",
		Value: url,
	})

	resourcespec := pipeline.PipelineResourceSpec{
		Type:   "git",
		Params: inputparams,
	}
	resource.Spec = resourcespec
	return &resource
}

// createBuildTask create a task to clone a git repo and use Kaniko to build the docker image
func createBuildTask(image Image) pipeline.Task {
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
				Name: "workspace",
				Type: pipeline.PipelineResourceTypeGit,
			}},
			Params: []pipeline.TaskParam{
				{
					Name:    "pathToDockerFile",
					Default: "Dockerfile",
				},
				{
					Name: "pathToContext",
				},
			},
		},
		Outputs: &pipeline.Outputs{
			Resources: []pipeline.TaskResource{{
				Name: "image",
				Type: pipeline.PipelineResourceTypeImage,
			}},
		},
		Steps: []corev1.Container{{
			Name:    "build-and-push-" + convertName(image.BuildTaskName),
			Image:   "gcr.io/kaniko-project/executor",
			Command: []string{"/kaniko/executor"},
			Args: []string{
				"--dockerfile=${inputs.params.pathToDockerFile}",
				"--destination=${outputs.resources.image.url}",
				"--context=${inputs.params.pathToContext}",
			},
		}},
	}

	return task
}

func createContainer(task Task) corev1.Container {
	// Need to be a little more intelligent with the Image.
	path := task.Image.Path
	if task.Image.Type != DOCKER {
		path = registry + "/" + task.Image.PipelineResourceImage.Name
	}

	return corev1.Container{
		Name:    convertName(task.Identifier),
		Image:   path,
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
