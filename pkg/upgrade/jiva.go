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
	"os"
	"time"

	core "github.com/openebs/api/v2/pkg/kubernetes/core"
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	batchV1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

type jivaUpdateConfig struct {
	fromVersion        string
	toVersion          string
	namespace          string
	pvNames            []string
	backOffLimit       int32
	serviceAccountName string
	logLevel           int32
	additionalArgs     []string
}

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
	volNames, fromVersion, err := GetJivaVolumesVersion(k)
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
	cfg := jivaUpdateConfig{
		fromVersion:        fromVersion,
		toVersion:          upgradeOpts.ToVersion,
		namespace:          ns,
		pvNames:            volNames,
		serviceAccountName: "jiva-operator",
		backOffLimit:       4,
		logLevel:           4,
		additionalArgs:     addArgs(upgradeOpts),
	}

	// Check if a job is running with underlying PV
	res, err := checkIfJobIsAlreadyRunning(k, &cfg)
	// If error or upgrade job is already running return
	if err != nil || res {
		log.Fatal("An upgrade job is already running with the underlying volume!")
	}

	k.CreateBatchJob(BuildJivaBatchJob(&cfg), cfg.namespace)
}

// GetJivaVolumesVersion returns the Jiva volumes list and current version
func GetJivaVolumesVersion(k *client.K8sClient) ([]string, string, error) {
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

// Returns additional arguments like image-prefix and image-tags
func addArgs(upgradeOpts UpgradeOpts) []string {
	var result []string
	if upgradeOpts.ImagePrefix != "" {
		result = append(result, fmt.Sprintf("--to-version-image-prefix=%s", upgradeOpts.ImagePrefix))
	}

	if upgradeOpts.ImageTag != "" {
		result = append(result, fmt.Sprintf("--to-version-image-tag=%s", upgradeOpts.ImageTag))
	}

	return result
}

// BuildJivaBatchJob returns Job to be build
func BuildJivaBatchJob(cfg *jivaUpdateConfig) *batchV1.Job {
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

func getJivaContainerArguments(cfg *jivaUpdateConfig) []string {
	// Set container arguments
	args := append([]string{
		"jiva-volume",
		fmt.Sprintf("--from-version=%s", cfg.fromVersion),
		fmt.Sprintf("--to-version=%s", cfg.toVersion),
		"--v=4", // can be taken from flags
	}, cfg.pvNames...)
	args = append(args, cfg.additionalArgs...)
	return args
}

func checkIfJobIsAlreadyRunning(k *client.K8sClient, cfg *jivaUpdateConfig) (bool, error) {
	jobs, err := k.GetBatchJobs("", "cas-type=jiva,name=jiva-upgrade")
	if err != nil {
		return false, err
	}

	// runningJob holds the information about the jobs that are in use by the PV
	// that has an upgrade-job progress(any status) already going
	var runningJob *batchV1.Job
	func() {
		for _, job := range jobs.Items { // JobItems
			for _, pvName := range cfg.pvNames { // running pvs in control plane
				for _, container := range job.Spec.Template.Spec.Containers { // iterate on containers provided by the cfg
					for _, args := range container.Args { // check if the running jobs (PVs) and the upcoming job(PVs) are common
						if args == pvName {
							runningJob = &job
							return
						}
					}
				}
			}
		}
	}()

	if runningJob != nil {
		jobCondition := runningJob.Status.Conditions
		info := jobInfo{name: runningJob.Name, namespace: runningJob.Namespace}
		if runningJob.Status.Failed > 0 ||
			len(jobCondition) > 0 && jobCondition[0].Type == "Failed" && jobCondition[0].Status == "True" {
			fmt.Println("Previous job failed.")
			fmt.Println("Reason: ", getReason(runningJob))
			fmt.Println("Creating a new Job with name:", info.name)
			// Job found thus delete the job and return false so that further process can be started
			if err := startDeletionTask(k, &info); err != nil {
				fmt.Println("error deleting job:", err)
				return true, err
			}
		}

		if runningJob.Status.Active > 0 {
			fmt.Println("A job is already active with the name ", runningJob.Name, " that is upgrading the PV.")
			// TODO:  Check the POD underlying the PV if their is any error inside
			return true, nil
		}

		if runningJob.Status.Succeeded > 0 {
			fmt.Println("Previous upgrade-job was successful for upgrading P.V.")
			// Provide the option to restart the Job
			shouldStart := util.PromptToStartAgain("Do you want to restart the Job?(no)", false)
			if shouldStart {
				// Delete previous successful task
				if err := startDeletionTask(k, &info); err != nil {
					return true, err
				}
			} else {
				os.Exit(0)
			}
		}
		return false, nil
	}

	return false, nil
}

func getReason(job *batchV1.Job) string {
	reason := job.Status.Conditions[0].Reason
	if len(reason) == 0 {
		return "Reason Not Found, check by inspecting jobs"
	}
	return reason
}

// startDeletionTask instantiates a deletion process
func startDeletionTask(k *client.K8sClient, info *jobInfo) error {
	err := k.DeleteBatchJob(info.name, info.namespace)
	if err != nil {
		return err
	}
	confirmDeletion(k, info)
	return nil
}

// confirmDeletion runs until the job is successfully done or reached threshhold duration
func confirmDeletion(k *client.K8sClient, info *jobInfo) {
	// create interval to call function periodically
	interval := time.NewTicker(time.Second * 2)

	// Create channel
	channel := make(chan bool)

	// Set threshhold time
	go func() {
		time.Sleep(time.Second * 10)
		channel <- true
	}()

	for {
		select {
		case <-interval.C:
			_, err := k.GetBatchJob(info.name, info.namespace)
			// Job is deleted successfully
			if err != nil {
				return
			}
		case <-channel:
			fmt.Println("Waiting time reached! Try Again!")
			return
		}
	}
}
