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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cstortypes "github.com/openebs/api/v2/pkg/apis/types"
	"github.com/openebs/openebsctl/client"
	"github.com/openebs/openebsctl/kubectl-openebs/cli/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/printers"
)

var (
	pvcInfoCommandHelpText = `
This command fetches information and status  of  the  various  aspects 
of  the  PersistentVolumeClaims  and  its underlying related resources 
in the provided namespace. If no namespace is provided it uses default
namespace for execution.

$ kubectl openebs describe pvc [name1] [name2] ... [nameN] -n [namespace]

`
)

const (
	cstorPvcInfoTemplate = `
{{.Name}} Details :
------------------
NAME             : {{.Name}}
NAMESPACE        : {{.Namespace}}
CAS TYPE         : {{.CasType}}
BOUND VOLUME     : {{.BoundVolume}}
ATTACHED TO NODE : {{.AttachedToNode}}
POOL             : {{.Pool}}
STORAGE CLASS    : {{.StorageClassName}}
SIZE             : {{.Size}}
USED             : {{.Used}}
PV STATUS	 : {{.PVStatus}}

`

	detailsFromCVC = `
Additional Details from CVC :
-----------------------------
NAME          : {{ .metadata.name }}
REPLICA COUNT : {{ .spec.provision.replicaCount }}
POOL INFO     : {{ .status.poolInfo}}
VERSION       : {{ .versionDetails.status.current}}
UPGRADING     : {{if eq .versionDetails.status.current .versionDetails.desired}}false{{else}}true{{end}}
`

	genericPvcInfoTemplate = `
{{.Name}} Details :
------------------
NAME             : {{.Name}}
NAMESPACE        : {{.Namespace}}
CAS TYPE         : {{.CasType}}
BOUND VOLUME     : {{.BoundVolume}}
STORAGE CLASS    : {{.StorageClassName}}
SIZE             : {{.Size}}
PV STATUS    	 : {{.PVStatus}}
`
)

// NewCmdDescribePVC Displays the pvc describe details
func NewCmdDescribePVC() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pvc",
		Aliases: []string{"pvcs", "persistentvolumeclaims", "persistentvolumeclaim"},
		Short:   "Displays PersistentVolumeClaim information",
		Long:    pvcInfoCommandHelpText,
		Example: `kubectl openebs describe pvc cstor-vol-1 cstor-vol-2 -n storage`,
		Run: func(cmd *cobra.Command, args []string) {
			var pvNs, openebsNamespace string
			if pvNs, _ = cmd.Flags().GetString("namespace"); pvNs == "" {
				pvNs = "default"
			}
			openebsNamespace, _ = cmd.Flags().GetString("openebs-namespace")
			util.CheckErr(RunPVCInfo(cmd, args, pvNs, openebsNamespace), util.Fatal)
		},
	}
	return cmd
}

// RunPVCInfo runs info command and make call to display the results
func RunPVCInfo(cmd *cobra.Command, pvcs []string, ns string, openebsNs string) error {
	if len(pvcs) == 0 {
		return errors.New("Please give at least one pvc name to describe")
	}
	// TODO: Make K8sClient to be able to determine openebs namespace by itself.
	// Below is currently hardcoded to support only if resources are in openebs ns
	// because the -n flag is used to take the pvc namespace and same cannot be used to
	// take the openebs namespace
	clientset, err := client.NewK8sClient(openebsNs)
	util.CheckErr(err, util.Fatal)

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
			casType := util.GetCasTypeFromSC(sc)
			if casType == util.CstorCasType {
				if openebsNs == "" {
					nsFromCli, err := clientset.GetOpenEBSNamespace(util.CstorCasType)
					if err != nil {
						return errors.Wrap(err, "Error determining the openebs namespace, please specify using \"--openebs-namespace\" flag")
					}
					clientset.Ns = nsFromCli
				}
				// Create Empty template objects and fill gradually when underlying sub CRs are identified.
				pvcInfo := util.CstorPVCInfo{}

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
					pvcInfo.Size = util.ConvertToIBytes(cv.Spec.Capacity.String())
					pvcInfo.PVStatus = cv.Status.Phase
				}

				// fetching the underlying CStorVolumeConfig for the PV, to get the cvc info and Pool Name and notify the user
				// if the CStorVolumeConfig is not found.
				cvc, err := clientset.GetCVC(item.Spec.VolumeName)
				if err != nil {
					fmt.Println("Underlying CstorVolumeConfig is not found for: ", item.Name)
				} else {
					pvcInfo.Pool = cvc.Labels[cstortypes.CStorPoolClusterLabelKey]
				}

				// fetching the underlying CStorVolumeAttachment for the PV, to get the attached to node and notify the user
				// if the CStorVolumeAttachment is not found.
				cva, err := clientset.GetCVA(item.Spec.VolumeName)
				if err != nil {
					pvcInfo.AttachedToNode = "N/A"
					fmt.Println("Underlying CstorVolumeAttachment is not found for: ", item.Name)
				} else {
					pvcInfo.AttachedToNode = cva.Spec.Volume.OwnerNodeID
				}

				// fetching the underlying CStorVolumeReplicas for the PV, to list their details and notify the user
				// none of the replicas are running if the CStorVolumeReplicas are not found.
				cvrs, err := clientset.GetCVR(item.Spec.VolumeName)
				if err == nil && len(cvrs.Items) > 0 {
					pvcInfo.Used = util.ConvertToIBytes(util.GetUsedCapacityFromCVR(cvrs))
				}

				// Printing the Filled Details of the Cstor PVC
				err = util.PrintByTemplate("pvc", cstorPvcInfoTemplate, pvcInfo)
				if err != nil {
					return err
				}

				// fetching the underlying TargetPod for the PV, to display its relevant details and notify the user
				// if the TargetPod is not found.
				targetPod, err := clientset.GetCstorVolumeTargetPod(item.Name, item.Spec.VolumeName)
				if err == nil {
					fmt.Printf("Target Details :\n----------------\n")
					var rows []metav1.TableRow
					rows = append(rows, metav1.TableRow{Cells: []interface{}{targetPod.Namespace, targetPod.Name, util.GetReadyContainers(targetPod.Status.ContainerStatuses), targetPod.Status.Phase, util.Duration(time.Since(targetPod.ObjectMeta.CreationTimestamp.Time)), targetPod.Status.PodIP, targetPod.Spec.NodeName}})
					util.TablePrinter(util.CstorTargetDetailsColumnDefinations, rows, printers.PrintOptions{Wide: true})
				} else {
					fmt.Printf("Target Details :\n----------------\nNo target pod exists for the CstorVolume\n")
				}

				// If CVRs are found list them and show relevant details else notify the user none of the replicas are
				// running if not found
				if cvrs != nil && len(cvrs.Items) > 0 {
					fmt.Printf("\nReplica Details :\n-----------------\n")
					var rows []metav1.TableRow
					for _, cvr := range cvrs.Items {
						rows = append(rows, metav1.TableRow{Cells: []interface{}{cvr.Name, util.ConvertToIBytes(cvr.Status.Capacity.Total), util.ConvertToIBytes(cvr.Status.Capacity.Used), cvr.Status.Phase, util.Duration(time.Since(cvr.ObjectMeta.CreationTimestamp.Time))}})
					}
					util.TablePrinter(util.CstorReplicaColumnDefinations, rows, printers.PrintOptions{Wide: true})
				} else {
					fmt.Printf("\nReplica Details :\n-----------------\nNo running replicas found\n")
				}

				if cvc != nil {
					util.TemplatePrinter(detailsFromCVC, cvc)
				}

			} else {
				// TODO: Change below to support more casTypes.
				// Incase a non-cstor pvc is entered show minimal details pertaining to the PVC
				pvcInfo := util.PVCInfo{}
				pvcInfo.Name = item.Name
				pvcInfo.Namespace = item.Namespace
				pvcInfo.StorageClassName = *item.Spec.StorageClassName
				quantity := item.Status.Capacity[util.StorageKey]
				pvcInfo.Size = quantity.String()
				pv, err := clientset.GetPV(item.Spec.VolumeName)
				if err == nil {
					pvcInfo.BoundVolume = item.Spec.VolumeName
					pvcInfo.PVStatus = pv.Status.Phase
					pvcInfo.CasType = util.GetCasType(pv, sc)
				}
				err = util.PrintByTemplate("pvc", genericPvcInfoTemplate, pvcInfo)
				if err != nil {
					return err
				}
			}
			// A separator to separate multiple pvc describes
			if len(pvcList.Items) > 1 && ind != len(pvcList.Items)-1 {
				fmt.Println("-------------------------------------------------------------------------------------")
			}
		}
	} else {
		fmt.Println("No such PVCs were found in the", ns, "namespace")
	}
	return nil
}
