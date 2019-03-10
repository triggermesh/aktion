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

	v1alpha1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	tektonv1alpha1 "github.com/tektoncd/pipeline/pkg/client/clientset/versioned/typed/pipeline/v1alpha1"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rest "k8s.io/client-go/rest"
)

type TaskRunCreator struct{}

func (trc TaskRunCreator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	taskRefName := os.Getenv("TASK_NAME")
	namespace := os.Getenv("NAMESPACE")

	log.Infof("Start to create TaskRun with TaskName [%s] and namespace [%s]", taskRefName, namespace)

	c, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Created InClusterConfig: ", c)

	tekton, err := tektonv1alpha1.NewForConfig(c)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Created Tekton client: ", tekton)

	taskRuns := tekton.TaskRuns(namespace)

	tr := v1alpha1.TaskRun{
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
				Name: taskRefName,
			},
			Trigger: v1alpha1.TaskTrigger{
				Type: v1alpha1.TaskTriggerTypeManual,
				Name: "manual",
			},
		},
	}

	res, err := taskRuns.Create(&tr)
	if err != nil {
		log.Error(err)
	}

	log.Info("Created TaskRun object in Kubernetes API")

	out, err := yaml.Marshal(*res)
	if err != nil {
		log.Error(err)
	}

	log.Infof("TaskRun create output: %s", out)

	w.Write(out)
}

func main() {
	log.Info("Start server at port :8080 ")
	http.ListenAndServe(":8080", TaskRunCreator{})
}
