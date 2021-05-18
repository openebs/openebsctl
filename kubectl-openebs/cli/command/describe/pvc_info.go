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
Name          : {{.Name}}
Replica Count : {{.ReplicaCount}}
Pool Info     : {{.PoolInfo}}
Version       : {{.Version}}
Upgrading     : {{.Upgrading}}

`

	genericPvcInfoTemplate = `
{{.Name}} Details :
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
		Short:   "Displays PersistentVolumeClaim information",
		Long:    pvcInfoCommandHelpText,
		Example: `kubectl openebs describe pvc [name1] [name2] ... [nameN] -n [namespace]`,
		Run: func(cmd *cobra.Command, args []string) {
			var ns string // This namespace belongs to the PVC entered
			if ns, _ = cmd.Flags().GetString("namespace"); ns == "" {
				ns = "default"
			}
			util.CheckErr(RunPVCInfo(cmd, args, ns), util.Fatal)
		},
	}
	return cmd
}

func RunPVCInfo(cmd *cobra.Command, pvcs []string, ns string) error {
	// TODO: Make K8sClient to be able to determine openebs namespace by itself.
	// Below is currently hardcoded to support only if resources are in openebs ns
	// because the -n flag is used to take the pvc namespace and same cannot be used to
	// take the openebs namespace
	clientset, err := client.NewK8sClient("openebs")
	if err != nil {
		return errors.Wrap(err, "Failed to execute describe pvc command")
	}
	// Fetch the list of v1PVC objects by passing the name of PVCs taken through CLI in ns namespace
	pvcList, err := clientset.GetPVCs(ns, pvcs)
	// Incase the PVCs are not found no further operation to be performed
	if err != nil {
		return errors.Wrap(err, "Failed to execute describe pvc command")
	}
	if pvcList != nil && len(pvcList.Items) > 0 {
		// Loop over the PVCs that are valid and found in ns namespace and
		// show their details.
		for ind, item := range pvcList.Items {
			// Get the storageClass object and skip the rest of the loop if not found.
			sc, err := clientset.GetStorageClass(*item.Spec.StorageClassName)
			if err != nil {
				fmt.Println("Error Fetching StorageClass details for", item.Name)
				continue
			}
			// TODO: Adding support for other casTypes
			// Get the casType from the storage class and branch on basic of CSTOR and NON-CSTOR PVCs.
			casType := client.GetCasTypeFromSC(sc)
			if casType == util.CSTOR_CAS_TYPE {

				// Create Empty template objects and fill gradually when underlying sub CRs are identified.
				pvcInfo := util.CstorPVCInfo{}
				cvcInfo := util.CVCInfo{}

				pvcInfo.Name = item.Name
				pvcInfo.Namespace = item.Namespace
				pvcInfo.BoundVolume = item.Spec.VolumeName
				pvcInfo.CasType = casType
				pvcInfo.StorageClassName = *item.Spec.StorageClassName

				// fetching the underlying CStorVolume for the PV, to get the phase and size and notify the user
				// if the CStorVolume is not found.
				cv, err := clientset.GetcStorVolume(item.Spec.VolumeName)
				if err != nil {
					fmt.Println("Underlying CstorVolume is not found for: ", item.Name)
				} else {
					pvcInfo.Size = cv.Spec.Capacity.String()
					pvcInfo.PVStatus = cv.Status.Phase
				}

				// fetching the underlying CStorVolumeConfig for the PV, to get the cvc info and Pool Name and notify the user
				// if the CStorVolumeConfig is not found.
				cvc, err := clientset.GetCVC(item.Spec.VolumeName)
				if err != nil {
					fmt.Println("Underlying CstorVolumeConfig is not found for: ", item.Name)
				} else {
					pvcInfo.Pool = cvc.Labels[cstortypes.CStorPoolClusterLabelKey]
					cvcInfo.Name = cvc.Name
					cvcInfo.ReplicaCount = len(cvc.Status.PoolInfo)
					cvcInfo.PoolInfo = cvc.Status.PoolInfo
					cvcInfo.Version = cvc.VersionDetails.Status.Current
					cvcInfo.Upgrading = !(cvc.VersionDetails.Status.Current == cvc.VersionDetails.Desired)
				}

				// fetching the underlying CStorVolumeAttachment for the PV, to get the attached to node and notify the user
				// if the CStorVolumeAttachment is not found.
				cva, err := clientset.GetCVA(item.Spec.VolumeName)
				if err != nil {
					fmt.Println("Underlying CstorVolumeAttachment is not found for: ", item.Name)
				} else {
					pvcInfo.AttachedToNode = cva.Spec.Volume.OwnerNodeID
				}

				// fetching the underlying CStorVolumeReplicas for the PV, to list their details and notify the user
				// none of the replicas are running if the CStorVolumeReplicas are not found.
				cvrs, err := clientset.GetCVR(item.Spec.VolumeName)
				if err == nil && len(cvrs.Items) > 0 {
					pvcInfo.Used = client.GetUsedCapacityFromCVR(cvrs)
				}

				// Printing the Filled Details of the Cstor PVC
				err = util.PrintByTemplate("pvc", cstorPvcInfoTemplate, pvcInfo)
				if err != nil {
					return err
				}

				// fetching the underlying TargetPod for the PV, to display its relevant details and notify the user
				// if the TargetPod is not found.
				targetPod, err := clientset.GetCstorVolumeTargetPod(item.Spec.VolumeName)
				if err == nil {
					targetPodOutput := make([]string, 2)
					fmt.Printf("Target Details :\n")
					targetPodOutput[0] = "Namespace|Name|Ready|Status|Age|IP|Node"
					targetPodOutput[1] = fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s",
						targetPod.Namespace,
						targetPod.Name,
						client.GetReadyContainers(targetPod.Status.ContainerStatuses),
						targetPod.Status.Phase,
						util.Duration(time.Since(targetPod.ObjectMeta.CreationTimestamp.Time)),
						targetPod.Status.PodIP,
						targetPod.Spec.NodeName,
					)
					fmt.Println(util.FormatList(targetPodOutput))
				} else {
					fmt.Println("Target Details :\nNo target pod exists for the CstorVolume")
				}

				// If CVRs are found list them and show relevant details else notify the user none of the replicas are
				// running if not found
				if cvrs != nil && len(cvrs.Items) > 0 {
					fmt.Printf("\nReplica Details :\n")
					cvrOutput := make([]string, len(cvrs.Items)+1)
					cvrOutput[0] = "Name|Total|Used|Status|Age"
					for i, cvr := range cvrs.Items {
						cvrOutput[i+1] = fmt.Sprintf("%s|%s|%s|%s|%s",
							cvr.Name,
							cvr.Status.Capacity.Total,
							cvr.Status.Capacity.Used,
							cvr.Status.Phase,
							util.Duration(time.Since(cvr.ObjectMeta.CreationTimestamp.Time)),
						)
					}
					fmt.Println(util.FormatList(cvrOutput))
				} else {
					fmt.Println("\nReplica Details :\nNo running replicas found")
				}

				if cvc != nil {
					// Printing the Filled Details of the CstorVolumeConfig
					err = util.PrintByTemplate("cvc", detailsFromCVC, cvcInfo)
					if err != nil {
						return err
					}
				}

			} else {
				// TODO: Change below to support more casTypes.
				// Incase a non-cstor pvc is entered show minimal details pertaining to the PVC
				pvcInfo := util.PVCInfo{}
				pvcInfo.Name = item.Name
				pvcInfo.Namespace = item.Namespace
				pvcInfo.StorageClassName = *item.Spec.StorageClassName
				quantity := item.Status.Capacity[util.STORAGE]
				pvcInfo.Size = quantity.String()
				pv, err := clientset.GetPV(item.Spec.VolumeName)
				if err == nil {
					pvcInfo.BoundVolume = item.Spec.VolumeName
					pvcInfo.PVStatus = pv.Status.Phase
					pvcInfo.CasType = client.GetCasType(pv, sc)
				}
				err = util.PrintByTemplate("pvc", genericPvcInfoTemplate, pvcInfo)
				if err != nil {
					return err
				}
			}
			if len(pvcList.Items) > 1 && ind != len(pvcList.Items)-1 {
				fmt.Println("-------------------------------------------------------------------------------------")
			}
		}
	} else {
		fmt.Println("No such PVCs were found in the", ns, "namespace")
	}
	return nil
}
