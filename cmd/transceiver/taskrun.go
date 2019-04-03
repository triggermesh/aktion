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

package main

import (
	"net/http"
	"os"
	"time"

	//v1alpha1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	//tektonv1alpha1 "github.com/tektoncd/pipeline/pkg/client/clientset/versioned/typed/pipeline/v1alpha1"
	v1alpha1 "github.com/knative/build-pipeline/pkg/apis/pipeline/v1alpha1"
	tektonv1alpha1 "github.com/knative/build-pipeline/pkg/client/clientset/versioned/typed/pipeline/v1alpha1"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
)

//TaskRunCreator handles TaskRun objects creation
type TaskRunCreator struct{}

func main() {
	log.Info("Start server at port :8080 ")
	http.ListenAndServe(":8080", TaskRunCreator{})
}

func (trc TaskRunCreator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	out, err := createTaskRuns()
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("TaskRun create output: %s", out)
	w.Write(out)
}

func createTaskRuns() ([]byte, error) {
	namespace := os.Getenv("NAMESPACE")

	c, err := rest.InClusterConfig()
	if err != nil {
		return []byte{}, err
	}
	log.Info("Created InClusterConfig: ", c)

	tekton, err := tektonv1alpha1.NewForConfig(c)
	if err != nil {
		return []byte{}, err
	}

	log.Info("Created Tekton client: ", tekton)

	simpleTaskrun, taskruns := v1alpha1.TaskRun{}, []v1alpha1.TaskRun{}
	// if TASKRUN_CONFIGMAP env variable is set, try to parse its content into taskruns list
	if taskrunConfigmap, ok := os.LookupEnv("TASKRUN_CONFIGMAP"); ok {
		core, err := kubernetes.NewForConfig(c)
		if err != nil {
			return []byte{}, err
		}
		configmap, err := core.CoreV1().ConfigMaps(namespace).Get(taskrunConfigmap, metav1.GetOptions{})
		if err != nil {
			return []byte{}, err
		}
		taskruns, _ = taskrunsFromConfigmaps(namespace, configmap)
	}

	// if TASK_NAME env variable is set, generate simple taskrun with TASK_NAME referrence
	if taskRefName, ok := os.LookupEnv("TASK_NAME"); ok {
		simpleTaskrun = taskRunWithTaskRef(namespace, taskRefName)
	}
	taskruns = append(taskruns, simpleTaskrun)

	// iterate over taskruns list and create its items
	var result []*v1alpha1.TaskRun
	for _, taskrun := range taskruns {
		tr, err := tekton.TaskRuns(namespace).Create(&taskrun)
		if err != nil {
			return []byte{}, err
		}
		result = append(result, tr)
		log.Infof("Created %s TaskRun object in Kubernetes API\n", tr.Name)
	}

	return yaml.Marshal(result)
}

func taskrunsFromConfigmaps(namespace string, configmap *corev1.ConfigMap) ([]v1alpha1.TaskRun, error) {
	var res []v1alpha1.TaskRun
	for _, v := range configmap.Data {
		var taskrun v1alpha1.TaskRun
		if err := yaml.Unmarshal([]byte(v), &taskrun); err != nil {
			return res, err
		}
		res = append(res, taskrun)
	}
	return res, nil
}

func taskRunWithTaskRef(namespace string, taskRef string) v1alpha1.TaskRun {
	return v1alpha1.TaskRun{
		TypeMeta: metav1.TypeMeta{
			Kind:       "TaskRun",
			APIVersion: "pipeline.knative.dev/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName:      "task-run-",
			Namespace:         namespace,
			CreationTimestamp: metav1.Time{Time: time.Now()},
		},
		Spec: v1alpha1.TaskRunSpec{
			TaskRef: &v1alpha1.TaskRef{
				Name: taskRef,
			},
			Trigger: v1alpha1.TaskTrigger{
				Type: v1alpha1.TaskTriggerTypeManual,
				Name: "manual",
			},
		},
	}
}
