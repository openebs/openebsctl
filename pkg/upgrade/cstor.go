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
	"errors"
	"fmt"
	"log"

	cstorv1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	"github.com/openebs/api/v2/pkg/kubernetes/core"
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	batchV1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

func InstantiateCspcUpgrade(options UpgradeOpts) {
	k := client.NewK8sClient()

	// auto-determine cstor namespace
	var err error
	k.Ns, err = k.GetOpenEBSNamespace(util.CstorCasType)
	if err != nil {
		fmt.Println(`Error determining cstor namespace! using "openebs" as namespace`)
		k.Ns = "openebs"
	}

	cspcList, err := k.ListCSPC()
	if err != nil {
		log.Fatal("err listing CSPC ", err)
	}

	poolNames := getCSPCPoolNames(cspcList)
	cfg := UpgradeJobCfg{
		fromVersion:        "",
		toVersion:          "",
		namespace:          k.Ns,
		resources:          poolNames,
		serviceAccountName: "",
		backOffLimit:       4,
		logLevel:           4,
		additionalArgs:     addArgs(options),
	}

	cfg.fromVersion, cfg.toVersion, err = getCstorVersionDetails(cspcList)
	if err != nil {
		fmt.Println("error: ", err)
	}
	if options.ToVersion != "" { // overriding the desired version from the cli flag
		cfg.toVersion = options.ToVersion
	}

	operator, err := k.GetCSPCOperator()
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	cfg.serviceAccountName = operator.Spec.ServiceAccountName

	// Check if a job is running with underlying PV
	err = inspectRunningUpgradeJobs(k, &cfg)
	// If error or upgrade job is already running return
	if err != nil {
		log.Fatal("An upgrade job is already running with the underlying volume!, More: ", err)
	}

	// Create upgrade job
	k.CreateBatchJob(buildCspcbatchJob(&cfg), k.Ns)
}

// buildCspcbatchJob returns CSPC Job to be build
func buildCspcbatchJob(cfg *UpgradeJobCfg) *batchV1.Job {
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

func getCstorCspcContainerArgs(cfg *UpgradeJobCfg) []string {
	// Set container arguments
	args := append([]string{
		"cstor-cspc",
		fmt.Sprintf("--from-version=%s", cfg.fromVersion),
		fmt.Sprintf("--to-version=%s", cfg.toVersion),
		"--v=4", // can be taken from flags
	}, cfg.resources...)
	args = append(args, cfg.additionalArgs...)
	return args
}

func getCSPCPoolNames(cspcList *cstorv1.CStorPoolClusterList) []string {
	var poolNames []string
	for _, cspc := range cspcList.Items {
		poolNames = append(poolNames, cspc.Name)
	}

	return poolNames
}

// getCstorVersionDetails returns cstor versioning details for upgrade job cfg
// It returns fromVersion, toVersion, or error
func getCstorVersionDetails(cspcList *cstorv1.CStorPoolClusterList) (fromVersion string, toVersion string, err error) {
	fmt.Println("Fetching CSPC control plane and Data Plane Version")
	for _, cspc := range cspcList.Items {
		fromVersion = cspc.VersionDetails.Status.Current
		toVersion = cspc.VersionDetails.Desired

		if fromVersion != "" && toVersion != "" {
			fmt.Println("Current Version:", fromVersion)
			fmt.Println("Desired Version:", toVersion)
			return
		}
	}

	return "", "", errors.New("problems fetching versioning details")
}

func GetCSPCOperatorServiceAccName(k *client.K8sClient) string {
	pods, err := k.GetPods("openebs.io/component-name=cspc-operator", "", k.Ns)
	if err != nil || len(pods.Items) == 0 {
		log.Fatal("error occured while searching operator, or no operator is found: ", err)
	}

	return pods.Items[0].Spec.ServiceAccountName
}
