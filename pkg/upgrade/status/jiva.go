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

package status

import (
	"fmt"

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
)

// Get job with the name -> apply selector to pod
func GetJobStatus(namespace string) {
	k := client.NewK8sClient()
	k.Ns = namespace
	// get jiva-upgrade batch jobs
	joblist, err := k.GetBatchJobs(namespace, "cas-type=jiva,name=jiva-upgrade")
	if err != nil {
		fmt.Println("Error getting jiva-upgrade jobs:", err)
		return
	}

	// No jobs found
	if len(joblist.Items) == 0 {
		fmt.Printf("No upgrade-jobs Found in %s namespace", namespace)
		return
	}

	for _, job := range joblist.Items {
		fmt.Println("***************************************")
		fmt.Println("Job Name: ", job.Name)
		getPodLogs(k, job.Name, namespace)
	}
	fmt.Println("***************************************")
}

// Get all the logs from the pods associated with a job
func getPodLogs(k *client.K8sClient, name string, namespace string) {
	// get pods created by the job
	podList, err := k.GetPods(fmt.Sprintf("job-name=%s", name), "", namespace)
	if err != nil {
		printColoredText(fmt.Sprintf("error getting pods of job %s, err: %s", name, err), util.Red)
		return
	}

	// range over pods to get all the logs
	for _, pod := range podList.Items {
		fmt.Println("From Pod:", pod.Name)
		logs := k.GetPodLogs(pod.Name, namespace)
		if logs == "" {
			fmt.Printf("-> No recent logs from the pod")
			fmt.Println()
			continue
		}
		printColoredText(logs, util.Blue)
	}

	if len(podList.Items) == 0 {
		printColoredText("No pods are running for this job", util.Red)
	}
}

func printColoredText(message string, color util.Color) {
	fmt.Println(util.ColorText(message, color))
}
