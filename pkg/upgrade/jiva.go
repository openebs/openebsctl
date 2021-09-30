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
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	batchV1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type jivaUpdateConfig struct {
	name               string
	fromVersion        string
	toVersion          string
	namespace          string
	pvNames            []string
	backOffLimit       int32
	serviceAccountName string
	logLevel           int32
}

// Jiva Data-plane Upgrade Job instantiator
func InstantiateJivaUpgrade(openebsNs string, toVersion string, menifestFile string) {
	k, err := client.NewK8sClient("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating k8s client")
		return
	}

	// p, _ := k.GetPods("job-name=jiva-volume-upgrade", "", "openebs")
	// fmt.Println(p)

	// If manifest Files is provided, apply the file to create a new upgrade-job
	if menifestFile != "" {
		yamlFile, err := yamlToJobSpec(menifestFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in Job: %s", err)
		}
		k.CreateBatchJob(yamlFile, yamlFile.Namespace)
		return
	}

	volNames, fromVersion, err := GetJivaVolumes(k)
	if err != nil {
		fmt.Println(err)
		return
	}

	// assign to-version
	if toVersion == "" {
		pods, e := k.GetPods("name=jiva-operator", "", "")
		if e != nil {
			fmt.Println("Failed to get operator-version, err: ", e)
			return
		}

		if len(pods.Items) == 0 {
			fmt.Println("Jiva-operator is not running!")
			return
		}

		toVersion = pods.Items[0].Labels["openebs.io/version"]
	}

	// assign namespace
	if openebsNs == "" {
		fmt.Println(`No Namespace Provided, using "default" as a namespace`)
		openebsNs = "default"
	}

	// create configuration
	n := fmt.Sprintf("jiva-upgrade-job-%v", rand.Intn(100)) // TODO: Seed random Numbers
	cfg := jivaUpdateConfig{
		name:               n,
		fromVersion:        fromVersion,
		toVersion:          toVersion,
		namespace:          openebsNs,
		pvNames:            volNames,
		serviceAccountName: "jiva-operator",
		backOffLimit:       4,
		logLevel:           4,
	}

	jobSpec := GetJivaBatchJob(&cfg)

	// Check if a job is running with underlying PV
	res, err := CheckIfJobIsAlreadyRunning(k, &cfg)
	// If error or upgrade job is already running return
	if err != nil || res {
		log.Fatal("An upgrade job is already running with the underlying volume!")
	}

	k.CreateBatchJob(jobSpec, cfg.namespace)
}

// GetJivaVolumes returns the Jiva volumes list and current version
func GetJivaVolumes(k *client.K8sClient) ([]string, string, error) {
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

func yamlToJobSpec(filePath string) (*batchV1.Job, error) {
	job := batchV1.Job{}
	// Check if the filepath is a remote-url
	if strings.HasPrefix(filePath, "http") {
		res, err := http.Get(filePath)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		// unmarshal yaml file into struct
		err = yaml.Unmarshal(body, &job)
		if err != nil {
			return nil, err
		}
	} else {
		// A file path is given located on local-disk of host
		yamlFile, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		// unmarshal yaml file to structs
		err = yaml.Unmarshal(yamlFile, &job)
		if err != nil {
			return nil, err
		}
	}

	return &job, nil
}

func getLatestJivaVersion() (string, error) {
	url := "https://raw.githubusercontent.com/openebs/jiva-operator/develop/deploy/helm/charts/Chart.yaml"

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var respData map[string]interface{}
	err = yaml.Unmarshal(body, &respData)
	if err != nil {
		return "", err
	}

	jivaLatestVersion := respData["version"].(string)
	return jivaLatestVersion, nil
}

// GetJivaBatchJob returns the Jiva Batch Specifications
func GetJivaBatchJob(cfg *jivaUpdateConfig) *batchV1.Job {
	jobSpec := &batchV1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cfg.name,
			Namespace: cfg.namespace,
		},
		Spec: batchV1.JobSpec{
			BackoffLimit: &cfg.backOffLimit,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					ServiceAccountName: cfg.serviceAccountName,
					Containers:         getJivaUpgradeContainer(cfg),
					RestartPolicy:      corev1.RestartPolicyOnFailure,
				},
			},
		},
	}

	return jobSpec
}

// getJivaUpgradeContainer returns containers for the jiva-upgrade-job
func getJivaUpgradeContainer(cfg *jivaUpdateConfig) []corev1.Container {
	return []corev1.Container{
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
	}
}

func CheckIfJobIsAlreadyRunning(k *client.K8sClient, cfg *jivaUpdateConfig) (bool, error) {
	jobs, err := k.GetBatchJobs()
	if err != nil {
		return false, err
	}

	var runningJob *batchV1.Job
	runningJobFound := false

	for _, job := range jobs.Items { // JobItems
		for _, pvName := range cfg.pvNames { // running pvs in control plane
			if !runningJobFound && !reflect.DeepEqual(job.Spec.Template, corev1.PodTemplateSpec{}) && !reflect.DeepEqual(job.Spec.Template.Spec, corev1.PodSpec{}) && len(job.Spec.Template.Spec.Containers) > 0 {
				for _, container := range job.Spec.Template.Spec.Containers { // iterate on containers provided by the cfg
					for _, args := range container.Args { // check if the running jobs (PVs) and the upcoming job(PVs) are common
						if args == pvName {
							runningJob = &job
							runningJobFound = true
							break
						}
					}
				}
			}
		}
	}

	if runningJobFound {
		active := runningJob.Status.Active
		failed := runningJob.Status.Failed
		succeeded := runningJob.Status.Succeeded

		if failed > 0 {
			fmt.Println("Previous job failed. Creating a new Job with name ", cfg.name, "...")
			// Job found but delete the job and return false so that further process can be started
			// TODO: Add Goroutines to handle job deletion completion
			return false, k.DeleteBatchJob(cfg.name, cfg.namespace)
		}

		if active > 0 {
			fmt.Println("A job is already active with the name ", runningJob.Name, " that is upgrading the PV")
			// TODO:  Check the POD underlying the PV if their is any error inside
			return true, nil
		}

		if succeeded > 0 {
			fmt.Println("Previous upgrade-job was successful for upgrading P.V., Not running current one.")
			os.Exit(0)
			// TODO:  Provide the option to restart the Job
		}
		return false, nil
	}

	return false, nil
}
