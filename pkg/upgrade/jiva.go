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

	core "github.com/openebs/api/v2/pkg/kubernetes/core"
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	batchV1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

type jobInfo struct {
	name      string
	namespace string
}

// Jiva Data-plane Upgrade Job instantiator
func InstantiateJivaUpgrade(upgradeOpts UpgradeOpts) {
	k := client.NewK8sClient()

	// auto-determine jiva namespace
	ns, err := k.GetOpenEBSNamespace(util.JivaCasType)
	if err != nil {
		fmt.Println(`Error determining namespace! using "openebs" as namespace`)
		ns = "openebs"
	}

	// get running volumes from cluster
	volNames, fromVersion, err := getJivaVolumesVersion(k)
	if err != nil {
		fmt.Println(err)
		return
	}

	// assign to-version
	if upgradeOpts.ToVersion == "" {
		pods, e := k.GetPods("name=jiva-operator", "", "")
		if e != nil {
			fmt.Println("Failed to get operator-version, err: ", e)
			return
		}

		if len(pods.Items) == 0 {
			fmt.Println("Jiva-operator is not running!")
			return
		}

		upgradeOpts.ToVersion = pods.Items[0].Labels["openebs.io/version"]
	}

	// create configuration
	cfg := UpgradeJobCfg{
		fromVersion:        fromVersion,
		toVersion:          upgradeOpts.ToVersion,
		namespace:          ns,
		resources:          volNames,
		serviceAccountName: "jiva-operator",
		backOffLimit:       4,
		logLevel:           4,
		additionalArgs:     addArgs(upgradeOpts),
	}

	// Check if a job is running with underlying PV
	res, err := inspectRunningUpgradeJobs(k, &cfg)
	// If error or upgrade job is already running return
	if err != nil || res {
		log.Fatal("An upgrade job is already running with the underlying volume!")
	}

	k.CreateBatchJob(BuildJivaBatchJob(&cfg), cfg.namespace)
}

// getJivaVolumesVersion returns the Jiva volumes list and current version
func getJivaVolumesVersion(k *client.K8sClient) ([]string, string, error) {
	// 1. Fetch all jivavolumes CRs in all namespaces
	_, jvMap, err := k.GetJVs(nil, util.Map, "", util.MapOptions{Key: util.Name})
	if err != nil {
		return nil, "", fmt.Errorf("err getting jiva volumes: %s", err.Error())
	}

	var jivaList *corev1.PersistentVolumeList
	//2. Get Jiva Persistent volumes
	jivaList, err = k.GetPvByCasType([]string{"jiva"}, "")
	if err != nil {
		return nil, "", fmt.Errorf("err getting jiva volumes: %s", err.Error())
	}

	var volumeNames []string
	var version string

	//3. Write-out names, versions and desired-versions
	for _, pv := range jivaList.Items {
		volumeNames = append(volumeNames, pv.Name)
		if v, ok := jvMap[pv.Name]; ok && len(version) == 0 {
			version = v.VersionDetails.Status.Current
		}
	}

	//4. Check for zero jiva-volumes
	if len(version) == 0 || len(volumeNames) == 0 {
		return volumeNames, version, fmt.Errorf("no jiva volumes found")
	}

	return volumeNames, version, nil
}

// BuildJivaBatchJob returns Job to be build
func BuildJivaBatchJob(cfg *UpgradeJobCfg) *batchV1.Job {
	return NewJob().
		WithGeneratedName("jiva-upgrade").
		WithLabel(map[string]string{"name": "jiva-upgrade", "cas-type": "jiva"}). // sets labels for job discovery
		WithNamespace(cfg.namespace).
		WithBackOffLimit(cfg.backOffLimit).
		WithPodTemplateSpec(
			func() *core.PodTemplateSpec {
				return core.NewPodTemplateSpec().
					WithServiceAccountName(cfg.serviceAccountName).
					WithContainers(
						func() *core.Container {
							return core.NewContainer().
								WithName("upgrade-jiva-go").
								WithArgumentsNew(getJivaContainerArguments(cfg)).
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

func getJivaContainerArguments(cfg *UpgradeJobCfg) []string {
	// Set container arguments
	args := append([]string{
		"jiva-volume",
		fmt.Sprintf("--from-version=%s", cfg.fromVersion),
		fmt.Sprintf("--to-version=%s", cfg.toVersion),
		"--v=4", // can be taken from flags
	}, cfg.resources...)
	args = append(args, cfg.additionalArgs...)
	return args
}
