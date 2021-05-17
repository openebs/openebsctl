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

package describe

import (
	"fmt"
	"html/template"
	"os"
	"time"

	cstortypes "github.com/openebs/api/v2/pkg/apis/types"
	"github.com/openebs/openebsctl/client"
	"github.com/openebs/openebsctl/kubectl-openebs/cli/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	pvcInfoCommandHelpText = `
This command fetches information and status of the various
aspects of a PersistentVolumeClaim and its underlying related
resources.

#
$ kubectl openebs describe pvc [name1] [name2] ... [nameN] -n [namespace]

`
)

const (
	cstorPvcInfoTemplate = `
{{.Name}} Details :
-------------------
Name             : {{.Name}}
Namespace        : {{.Namespace}}
Cas Type         : {{.CasType}}
Bound Volume     : {{.BoundVolume}}
Attached To Node : {{.AttachedToNode}}
Pool             : {{.Pool}}
Storage Class    : {{.StorageClassName}}
Size             : {{.Size}}
Used             : {{.Used}}
PV Status	 : {{.PVStatus}}

`

	detailsFromCVC = `
Additional Details from CVC :
-----------------------------
Name          : {{.Name}}
Replica Count : {{.ReplicaCount}}
Pool Info     : {{.PoolInfo}}
Version       : {{.Version}}
Upgrading     : {{.Upgrading}}

`

	genericPvcInfoTemplate = `
{{.Name}} Details :
-------------------
Name             : {{.Name}}
Namespace        : {{.Namespace}}
Cas Type         : {{.CasType}}
Bound Volume     : {{.BoundVolume}}
Storage Class    : {{.StorageClassName}}
Size             : {{.Size}}
PV Status	 : {{.PVStatus}}

`
)

func NewCmdDescribePVC() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pvc",
		Aliases: []string{"pvcs", "persistentvolumeclaims", "persistentvolumeclaim"},
		Short:   "Displays Openebs information",
		Long:    volumeInfoCommandHelpText,
		Example: `kubectl openebs describe pvc [name1] [name2] ... [nameN] -n [namespace]`,
		Run: func(cmd *cobra.Command, args []string) {
			var ns string
			if ns, _ = cmd.Flags().GetString("namespace"); ns == "" {
				// NOTE: The error comes as nil even when the ns flag is not specified
				ns = "default"
			}
			util.CheckErr(RunPVCInfo(cmd, args, ns), util.Fatal)
		},
	}
	return cmd
}

func RunPVCInfo(cmd *cobra.Command, pvcs []string, ns string) error {
	// TODO: Make K8sClient to be able to determine openebs namespace by itself.
	// Please change the below with the namespace where openebs is installed in the cluster
	clientset, err := client.NewK8sClient("openebs")
	if err != nil {
		return errors.Wrap(err, "failed to execute pvc info command")
	}
	pvcList, err := clientset.GetPVCs(ns, pvcs)
	if err != nil {
		return errors.Wrap(err, "failed to execute pvc info command")
	}
	for _, item := range pvcList.Items {
		sc, err := clientset.GetStorageClass(*item.Spec.StorageClassName)
		if err != nil {
			fmt.Println("Error Fetching StorageClass details for", item.Name)
			continue
		}
		casType := client.GetCasTypeFromSC(sc)
		if casType == util.CSTOR_CAS_TYPE {
			pvcInfo := util.CstorPVCInfo{}
			cv, err := clientset.GetcStorVolume(item.Spec.VolumeName)
			if err != nil {
				fmt.Println("Error Fetching ctsor volume details for", item.Name)
				continue
			}

			cvcInfo := util.CVCInfo{}
			cvc, err := clientset.GetCVC(item.Spec.VolumeName)
			if err != nil {
				fmt.Println("Error Fetching cstor volume config details for", item.Name)
				continue
			}
			cvcInfo.Name = cvc.Name
			cvcInfo.ReplicaCount = len(cvc.Status.PoolInfo)
			cvcInfo.PoolInfo = cvc.Status.PoolInfo
			cvcInfo.Version = cvc.VersionDetails.Status.Current
			cvcInfo.Upgrading = !(cvc.VersionDetails.Status.Current == cvc.VersionDetails.Desired)

			cva, err := clientset.GetCVA(item.Spec.VolumeName)
			if err == nil {
				pvcInfo.AttachedToNode = cva.Spec.Volume.OwnerNodeID
			}
			cvrs, err := clientset.GetCVR(item.Spec.VolumeName)
			if err == nil && len(cvrs.Items) > 0 {
				pvcInfo.Used = client.GetUsedCapacityFromCVR(cvrs)
			}
			pvcInfo.Name = item.Name
			pvcInfo.Namespace = item.Namespace
			pvcInfo.BoundVolume = item.Spec.VolumeName
			pvcInfo.CasType = casType
			pvcInfo.Pool = cvc.Labels[cstortypes.CStorPoolClusterLabelKey]
			pvcInfo.AttachedToNode = cva.Spec.Volume.OwnerNodeID
			pvcInfo.StorageClassName = *item.Spec.StorageClassName
			pvcInfo.Size = cv.Spec.Capacity.String()
			pvcInfo.PVStatus = cv.Status.Phase

			pvcInfoTemplate, err := template.New("pvc").Parse(cstorPvcInfoTemplate)
			if err != nil {
				return errors.Wrap(err, "error displaying output for pvc info")
			}
			err = pvcInfoTemplate.Execute(os.Stdout, pvcInfo)
			if err != nil {
				return errors.Wrap(err, "error displaying cvc details")
			}

			targetPod, err := clientset.GetCstorVolumeTargetPod(item.Spec.VolumeName)
			targetPodOutput := make([]string, 3)
			if err == nil {
				fmt.Printf("Target Details :\n----------------\n")
				targetPodOutput[0] = "Namespace|Name|Ready|Status|Age|IP|Node"
				targetPodOutput[1] = "---------|----|-----|------|---|--|----"
				targetPodOutput[2] = fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s",
					targetPod.Namespace,
					targetPod.Name,
					client.GetReadyContainers(targetPod.Status.ContainerStatuses),
					targetPod.Status.Phase,
					util.Duration(time.Since(targetPod.ObjectMeta.CreationTimestamp.Time)),
					targetPod.Status.PodIP,
					targetPod.Spec.NodeName,
				)
			}
			fmt.Println(util.FormatList(targetPodOutput))

			if cvrs != nil {
				fmt.Printf("\nReplica Details :\n----------------\n")
				cvrOutput := make([]string, len(cvrs.Items)+2)
				cvrOutput[0] = "Name|Total|Used|Status|Age"
				cvrOutput[1] = "----|-----|----|------|---"
				for i, cvr := range cvrs.Items {
					cvrOutput[i+2] = fmt.Sprintf("%s|%s|%s|%s|%s",
						cvr.Name,
						cvr.Status.Capacity.Total,
						cvr.Status.Capacity.Used,
						cvr.Status.Phase,
						util.Duration(time.Since(cvr.ObjectMeta.CreationTimestamp.Time)),
					)
				}
				fmt.Println(util.FormatList(cvrOutput))
			}

			cvcInfoTemplate, err := template.New("cvc").Parse(detailsFromCVC)
			if err != nil {
				return errors.Wrap(err, "error displaying output for cvc info")
			}
			err = cvcInfoTemplate.Execute(os.Stdout, cvcInfo)
			if err != nil {
				return errors.Wrap(err, "error displaying cvc details")
			}
		} else {
			pvcInfo := util.PVCInfo{}
			pvcInfo.Name = item.Name
			pvcInfo.Namespace = item.Namespace
			pvcInfo.StorageClassName = *item.Spec.StorageClassName
			quantity := item.Status.Capacity["storage"]
			pvcInfo.Size = quantity.String()
			pv, err := clientset.GetPV(item.Spec.VolumeName)
			if err == nil {
				pvcInfo.BoundVolume = item.Spec.VolumeName
				pvcInfo.PVStatus = pv.Status.Phase
				pvcInfo.CasType = client.GetCasType(pv, sc)
			}
			cvcInfoTemplate, err := template.New("pvc").Parse(genericPvcInfoTemplate)
			if err != nil {
				return errors.Wrap(err, "error displaying output for pvc info")
			}
			err = cvcInfoTemplate.Execute(os.Stdout, pvcInfo)
			if err != nil {
				return errors.Wrap(err, "error displaying pvc details")
			}
		}
	}
	return nil
}
