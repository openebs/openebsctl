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
	"log"

	"github.com/openebs/api/v2/pkg/kubernetes/core"
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	batchV1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

type cstorUpdateConfig struct {
	fromVersion        string
	toVersion          string
	namespace          string
	poolNames          []string
	backOffLimit       int32
	serviceAccountName string
	logLevel           int32
	additionalArgs     []string
}

func InstantiateCspcUpgrade(options UpgradeOpts) {
	k := client.NewK8sClient()

	// auto-determine cstor namespace
	var err error
	k.Ns, err = k.GetOpenEBSNamespace(util.CstorCasType)
	if err != nil {
		fmt.Println(`Error determining cstor namespace! using "openebs" as namespace`)
		k.Ns = "openebs"
	}

	poolNames := getCSPCPoolNames(k)
	cfg := cstorUpdateConfig{
		fromVersion:        "",
		toVersion:          "",
		namespace:          k.Ns,
		poolNames:          poolNames,
		serviceAccountName: "openebs-maya-operator",
		backOffLimit:       4,
		logLevel:           4,
		additionalArgs:     addArgs(options),
	}

}

func getCSPCPoolNames(k *client.K8sClient) []string {
	scList, err := k.GetScWithCasType(util.CstorCasType)
	if err != nil {
		log.Fatal("err listing storage classes: ", err)
	}

	// Set to contain cspc names
	cspcNames := make(map[string]bool)
	for _, sc := range scList {
		cspcName := sc.Parameters["cstorPoolCluster"]
		if cspcName != "" {
			cspcNames[cspcName] = true
		}
	}

	// create slice and return it
	poolnames := make([]string, len(cspcNames))
	i := 0
	for pool := range cspcNames {
		poolnames[i] = pool
		i++
	}

	return poolnames
}

// buildCspcJobConfig returns CSPC Job to be build
func buildCspcJobConfig(cfg *cstorUpdateConfig) *batchV1.Job {
	return NewJob().
		WithGeneratedName("cstor-cspc-upgrade").
		WithLabel(map[string]string{"name": "cstor-cspc-upgrade", "cas-type": "cstor"}). // sets labels for job discovery
		WithNamespace(cfg.namespace).
		WithBackOffLimit(cfg.backOffLimit).
		WithPodTemplateSpec(
			func() *core.PodTemplateSpec {
				return core.NewPodTemplateSpec().
					WithServiceAccountName(cfg.serviceAccountName).
					WithContainers(
						func() *core.Container {
							return core.NewContainer().
								WithName("upgrade-cstor-cspc-go").
								WithArgumentsNew(getCstorCspcContainerArgs(cfg)).
								WithEnvsNew(
									[]corev1.EnvVar{
										{
											Name: "OPENEBS_NAMESPACE",
											ValueFrom: &corev1.EnvVarSource{
												FieldRef: &corev1.ObjectFieldSelector{
													FieldPath: "metadata.namespace",
												},
											},
										},
									},
								).
								WithImage(fmt.Sprintf("openebs/upgrade:%s", cfg.toVersion)).
								WithImagePullPolicy(corev1.PullIfNotPresent) // Add TTY to openebs/api
						}(),
					)
			}(),
		).
		WithRestartPolicy(corev1.RestartPolicyOnFailure). // Add restart policy in openebs/api
		Job
}

func getCstorCspcContainerArgs(cfg *cstorUpdateConfig) []string {
	// Set container arguments
	args := append([]string{
		"cstor-cspc",
		fmt.Sprintf("--from-version=%s", cfg.fromVersion),
		fmt.Sprintf("--to-version=%s", cfg.toVersion),
		"--v=4", // can be taken from flags
	}, cfg.poolNames...)
	args = append(args, cfg.additionalArgs...)
	return args
}
