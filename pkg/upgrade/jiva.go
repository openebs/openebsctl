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
	"time"

	"github.com/manifoldco/promptui"
	core "github.com/openebs/api/v2/pkg/kubernetes/core"
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	batchV1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
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
	additionalArgs     []string
}

// The variables that is to be used by the CLI for jiva-upgrade-job
var CasType, ToVersion, File, ImagePrefix, ImageTag string

// Jiva Data-plane Upgrade Job instantiator
func InstantiateJivaUpgrade(openebsNs string) {
	k, err := client.NewK8sClient("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating k8s client")
		return
	}

	fmt.Println(CasType, ToVersion, File, ImagePrefix, ImageTag)

	// If manifest Files is provided, apply the file to create a new upgrade-job
	if File != "" {
		yamlFile, err := yamlToJobSpec(File)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in Job: %s", err)
		}
		k.CreateBatchJob(yamlFile, yamlFile.Namespace)
		return
	}

	// get running volumes from cluster
	volNames, fromVersion, err := GetJivaVolumes(k)
	if err != nil {
		fmt.Println(err)
		return
	}

	// assign to-version
	if ToVersion == "" {
		pods, e := k.GetPods("name=jiva-operator", "", "")
		if e != nil {
			fmt.Println("Failed to get operator-version, err: ", e)
			return
		}

		if len(pods.Items) == 0 {
			fmt.Println("Jiva-operator is not running!")
			return
		}

		ToVersion = pods.Items[0].Labels["openebs.io/version"]
	}

	// assign namespace
	if openebsNs == "" {
		fmt.Println(`No Namespace Provided, using "default" as a namespace`)
		openebsNs = "default"
	}

	// create configuration
	name := fmt.Sprintf("jiva-upgrade-job-%v", rand.Intn(100)) // TODO: Seed random Numbers
	cfg := jivaUpdateConfig{
		name:               name,
		fromVersion:        fromVersion,
		toVersion:          ToVersion,
		namespace:          openebsNs,
		pvNames:            volNames,
		serviceAccountName: "jiva-operator",
		backOffLimit:       4,
		logLevel:           4,
		additionalArgs:     addArgs(),
	}

	// Check if a job is running with underlying PV
	res, err := checkIfJobIsAlreadyRunning(k, &cfg)
	// If error or upgrade job is already running return
	if err != nil || res {
		log.Fatal("An upgrade job is already running with the underlying volume!")
	}

	k.CreateBatchJob(BuildJivaBatchJob(&cfg), cfg.namespace)
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

// Returns additional arguments like image-prefix and image-tags
func addArgs() []string {
	var result []string
	if ImagePrefix != "" {
		result = append(result, fmt.Sprintf("--to-version-image-prefix=%s", ImagePrefix))
	}

	if ImageTag != "" {
		result = append(result, fmt.Sprintf("--to-version-image-tag=%s", ImageTag))
	}

	return result
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

// BuildJivaBatchJob returns Job to be build
func BuildJivaBatchJob(cfg *jivaUpdateConfig) *batchV1.Job {
	return NewJob().
		WithName(cfg.name).
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
								WithArgumentsNew(getContainerArguments(cfg)).
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

func getContainerArguments(cfg *jivaUpdateConfig) []string {
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
			fmt.Println("Previous job failed. Creating a new Job with name:", cfg.name)
			// Job found thus delete the job and return false so that further process can be started
			if err := startDeletionTask(k, cfg); err != nil {
				return true, err
			}
		}

		if active > 0 {
			fmt.Println("A job is already active with the name ", runningJob.Name, " that is upgrading the PV")
			// TODO:  Check the POD underlying the PV if their is any error inside
			return true, nil
		}

		if succeeded > 0 {
			fmt.Println("Previous upgrade-job was successful for upgrading P.V.")
			// Provide the option to restart the Job
			shouldStart := promptToStartAgain()
			if shouldStart {
				// Delete previous successful task
				if err := startDeletionTask(k, cfg); err != nil {
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

// startDeletionTask instantiates a deletion process
func startDeletionTask(k *client.K8sClient, cfg *jivaUpdateConfig) error {
	err := k.DeleteBatchJob(cfg.name, cfg.namespace)
	if err != nil {
		return err
	}
	confirmDeletion(k, cfg)
	return nil
}

// confirmDeletion runs until the job is successfully done or reached threshhold duration
func confirmDeletion(k *client.K8sClient, cfg *jivaUpdateConfig) {
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
			_, err := k.GetBatchJob(cfg.name, cfg.namespace)
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

// promptToStartAgain returns if the job should be started Again
func promptToStartAgain() bool {
	prompt := promptui.Prompt{
		Label: "Do you want to restart the Job?(no)",
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Println("prompt Failed: ", err)
		return false
	}

	result = strings.ToLower(result)
	if result == "yes" || result == "y" {
		return true
	}

	return false
}
