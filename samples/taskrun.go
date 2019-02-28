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
	"context"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ghodss/yaml"
	"github.com/knative/build-pipeline/pkg/apis/pipeline/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//Handler handles events from GpcPubSub source
func Handler(ctx context.Context) ([]byte, error) {

	tr := v1alpha1.TaskRun{
		TypeMeta: metav1.TypeMeta{
			Kind:       "TaskRun",
			APIVersion: "pipeline.knative.dev/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              "taskname",
			Namespace:         "default",
			CreationTimestamp: metav1.Time{Time: time.Now()},
		},
		Spec: v1alpha1.TaskRunSpec{
			TaskRef: &v1alpha1.TaskRef{
				Name: "taskrefname",
			},
			Trigger: v1alpha1.TaskTrigger{
				Type: v1alpha1.TaskTriggerTypePipelineRun,
				Name: "testtriggername",
			},
		},
	}

	out, err := yaml.Marshal(tr)
	if err != nil {
		log.Fatal(err)
	}
	return out, nil
}

func main() {
	lambda.Start(Handler)
}
