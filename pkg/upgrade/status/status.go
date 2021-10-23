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
	"github.com/openebs/openebsctl/pkg/upgrade"
)

var WaitFlag bool // For opening wait stream for logs

// Get job with the name -> apply selector to pod
func GetJobStatus() {
	k := client.NewK8sClient()
	namespace := upgrade.OpenebsNs	

	// get jiva-upgrade batch jobs
	joblist, err := k.GetBatchJobs(namespace, "cas-type=jiva,name=jiva-upgrade")
	if err != nil {
		fmt.Println("Error getting jiva-upgrade jobs:", err)
		return
	}

	for _, job := range joblist.Items {
		fmt.Println("***************************************")
		fmt.Println("Job Name: ", job.Name)
		getPodLogs(k, job.Name, namespace)
		fmt.Println()
	}
	fmt.Println("***************************************")
}

// Get all the logs from the pods associated with a job
func getPodLogs(k *client.K8sClient, name string, namespace string) {
	// get pods created by the job
	podlist, err := k.GetPods(fmt.Sprintf("job-name=%s", name), "", namespace)
	if err != nil {
		fmt.Println("error getting pods of job", name, ": err", err)
		return
	}

	// range over pods to get all the logs
	for _, pod := range podlist.Items {
		fmt.Println("From Pod:", pod.Name)
		logs := k.GetPodLogs(pod, namespace)
		fmt.Println(logs)
	}
}
