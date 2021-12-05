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
	"time"

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	batchV1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

// UpgradeOpts are the upgrade options that are provided
// with the CLI flags
type UpgradeOpts struct {
	CasType     string
	ToVersion   string
	ImagePrefix string
	ImageTag    string
	ServiceAccountName string
}

// UpgradeJobCfg holds upgrade job confiogurations while creating a new Job
type UpgradeJobCfg struct {
	fromVersion        string
	toVersion          string
	namespace          string
	resources          []string
	backOffLimit       int32
	serviceAccountName string
	logLevel           int32
	additionalArgs     []string
}

// inspectRunningUpgradeJobs inspects all the jobs running in the cluster
// and returns if the job updating the resource is already available
func inspectRunningUpgradeJobs(k *client.K8sClient, cfg *UpgradeJobCfg) (bool, error) {
	jobs, err := k.GetBatchJobs("", "")
	if err != nil {
		return false, err
	}

	// runningJob holds the information about the jobs that are in use by the PV
	// that has an upgrade-job progress(any status) already going
	var runningJob *batchV1.Job
	func() {
		for _, job := range jobs.Items { // JobItems
			for _, pvName := range cfg.resources { // running pvs in control plane
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

	return runningJobHandler(k, runningJob)
}

// runningJobHandler checks the status of the job and takes action on it
// to modify or delete it based on the status of the Job
func runningJobHandler(k *client.K8sClient, runningJob *batchV1.Job) (bool, error) {

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

// getReason returns the reason for the current status of Job
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

// getServiceAccountName returns service account Name for the openEBS resource
func getServiceAccountName(podList *corev1.PodList) string {
	var serviceAccountName string
	for _, pod := range podList.Items {
		svname := pod.Spec.ServiceAccountName
		if svname != "" {
			serviceAccountName = svname
		}
	}
	return serviceAccountName
}
