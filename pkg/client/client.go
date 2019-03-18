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

package client

import (
	"os"

	pipelineApi "github.com/knative/build-pipeline/pkg/client/clientset/versioned"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	// gcp package is required for kubectl configs with GCP auth providers
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

//ConfigSet contains configurations from available configuration file or from in-cluster environment
type ConfigSet struct {
	Core     *kubernetes.Clientset
	Pipeline *pipelineApi.Clientset

	Config *rest.Config
}

//ConfigPath returns path to a valid config file
func ConfigPath(cfgFile string) string {
	homeDir := "."
	if dir := os.Getenv("HOME"); dir != "" {
		homeDir = dir
	}

	kubeconfig := os.Getenv("KUBECONFIG")
	if len(cfgFile) != 0 {
		// using config file passed with --config argument
	} else if _, err := os.Stat(kubeconfig); err == nil {
		cfgFile = kubeconfig
	} else {
		cfgFile = homeDir + "/.kube/config"
	}
	return cfgFile
}

// NewClient returns ConfigSet created from available configuration file or from in-cluster environment
func NewClient(cfgFile string) (ConfigSet, error) {
	config, err := clientcmd.BuildConfigFromFlags("", cfgFile)

	c := ConfigSet{
		Config: config,
	}

	if c.Pipeline, err = pipelineApi.NewForConfig(config); err != nil {
		return c, err
	}

	if c.Core, err = kubernetes.NewForConfig(config); err != nil {
		return c, err
	}

	return c, nil
}
