/*
Copyright 2020-2021 The OpenEBS Authors

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

package upgrade

import (
	"fmt"
	"os"

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/spf13/cobra"
	batchV1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type jivaUpdateConfig struct {
	fromVersion        string
	toVersion          string
	namespace          string
	pvNames            []string
	backOffLimit       int32
	serviceAccountName string
	logLevel           int32
}

// Jiva Data-plane Upgrade Job instantiator
func InstantiateJivaUpgrade(cmd *cobra.Command) {
	k, err := client.NewK8sClient("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating k8s client")
	}
	ns, err := cmd.Flags().GetString("openebs-namespace")
	handleErr(ns, "openebs-namespace", err)

	tv, err := cmd.Flags().GetString("to-version")
	handleErr(tv, "to-version", err)

	cType, err := cmd.Flags().GetString("cas-type")
	handleErr(cType, "cas-type", err)

	cfg := jivaUpdateConfig{
		fromVersion:        "2.7.0", // determine from control plane
		toVersion:          tv,
		namespace:          ns,
		pvNames:            []string{"pvc-9cebb2c3-b26e-4372-9e25-d1dc2d26c650"}, // determine from Control plane
		serviceAccountName: "jiva-operator",
		backOffLimit:       4,
		logLevel:           4,
	}

	jobSpec := GetJivaBatchJob(&cfg)
	k.CreateBatchJob(jobSpec)
}

// GetJivaBatchJob returns the Jiva Batch Specifications
func GetJivaBatchJob(cfg *jivaUpdateConfig) *batchV1.Job {
	var backOffLimit int32 = 5

	jobSpec := &batchV1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "jiva-volume-upgrade",
			Namespace: cfg.namespace,
		},
		Spec: batchV1.JobSpec{
			BackoffLimit: &backOffLimit,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					ServiceAccountName: cfg.serviceAccountName,
					Containers: []corev1.Container{
						{
							Name: "upgrade-jiva-go",
							Args: append([]string{
								"jiva-volume",
								fmt.Sprintf("--from-version=%s", cfg.fromVersion),
								fmt.Sprintf("--to-version=%s", cfg.toVersion),
								"--v=4", // can be taken from flags
							}, cfg.pvNames...),
							Env: []corev1.EnvVar{
								{
									Name: "OPENEBS_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
							},
							TTY:             true,
							Image:           fmt.Sprintf("openebs/upgrade:%s", cfg.toVersion),
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
					RestartPolicy: corev1.RestartPolicyOnFailure,
				},
			},
		},
	}

	return jobSpec
}

func handleErr(res string, argName string, err error) {
	if err != nil || len(res) == 0 {
		fmt.Fprintf(os.Stderr, "--%s not Provided or failed fetching --%s\n", argName, argName)
		fmt.Printf(`Try setting %s with "--%s" flag`, argName, argName)
		os.Exit(1)
		return
	}
}
