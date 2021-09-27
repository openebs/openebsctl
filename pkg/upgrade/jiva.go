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
	"net/http"
	"os"
	"strings"

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	batchV1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
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
func InstantiateJivaUpgrade(openebsNs string, toVersion string, menifestFile string) {
	k, err := client.NewK8sClient("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating k8s client")
		return
	}

	p, _ := k.GetPods("job-name=jiva-volume-upgrade", "", "openebs")
	fmt.Println(p)

	// If manifest Files is provided, apply the file to create a new upgrade-job
	if menifestFile != "" {
		yamlFile, err := yamlToJobSpec(menifestFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in Job: %s", err)
		}
		k.CreateBatchJob(yamlFile)
		return
	}

	volNames, fromVersion, desiredVersion, err := GetJivaVolumes(k)
	if err != nil {
		fmt.Println(err)
		return
	}

	if toVersion == "" {
		if desiredVersion != fromVersion {
			// Mark it as toVersion
			toVersion = desiredVersion
		} else {
			// TODO: Upgrade version to latest available version for Jiva volumes
			fmt.Println("Fetching latest version from the remote...")
			latVer, err := getLatestJivaVersion()
			if err != nil {
				log.Fatal("Error fetching latest version: ", err)
			}
			toVersion = latVer
		}
	}

	if openebsNs == "" {
		fmt.Println(`No Namespace Provided, using "default" as a namespace`)
		openebsNs = "default"
	}

	cfg := jivaUpdateConfig{
		fromVersion:        fromVersion,
		toVersion:          toVersion,
		namespace:          openebsNs,
		pvNames:            volNames,
		serviceAccountName: "jiva-operator",
		backOffLimit:       4,
		logLevel:           4,
	}

	jobSpec := GetJivaBatchJob(&cfg)
	k.CreateBatchJob(jobSpec)
}

// GetJivaVolumes returns the Jiva volumes list and current version
func GetJivaVolumes(k *client.K8sClient) ([]string, string, string, error) {
	// 1. Fetch all jivavolumes CRs in all namespaces
	_, jvMap, err := k.GetJVs(nil, util.Map, "", util.MapOptions{Key: util.Name})
	if err != nil {
		return nil, "", "", fmt.Errorf("err getting jiva volumes: %s", err.Error())
	}

	var jivaList *corev1.PersistentVolumeList
	//2. Get Jiva Persistent volumes
	jivaList, err = k.GetPvByCasType([]string{"jiva"}, "")
	if err != nil {
		return nil, "", "", fmt.Errorf("err getting jiva volumes: %s", err.Error())
	}

	var volumeNames []string
	var version, desiredVersion string

	//3. Write-out names, versions and desired-versions
	for _, pv := range jivaList.Items {
		volumeNames = append(volumeNames, pv.Name)
		if v, ok := jvMap[pv.Name]; ok && len(version) == 0 {
			version = v.VersionDetails.Status.Current
			desiredVersion = v.VersionDetails.Desired
		}
	}

	//4. Check for zero jiva-volumes
	if len(version) == 0 || len(volumeNames) == 0 {
		return volumeNames, version, desiredVersion, fmt.Errorf("no jiva volumes found")
	}

	return volumeNames, version, desiredVersion, nil
}

// GetJivaBatchJob returns the Jiva Batch Specifications
func GetJivaBatchJob(cfg *jivaUpdateConfig) *batchV1.Job {
	jobSpec := &batchV1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "jiva-volume-upgrade",
			Namespace: cfg.namespace,
		},
		Spec: batchV1.JobSpec{
			BackoffLimit: &cfg.backOffLimit,
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
