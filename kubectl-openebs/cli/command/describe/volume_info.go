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
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"

	cstortypes "github.com/openebs/api/v2/pkg/apis/types"

	"github.com/openebs/openebsctl/client"
	"github.com/openebs/openebsctl/kubectl-openebs/cli/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	volumeInfoCommandHelpText = `
This command fetches information and status of the various
aspects of a cStor Volume such as ISCSI, Controller, and Replica.

#
$ kubectl openebs describe [volume] [names...]

`
)

const (
	volInfoTemplate = `
{{.Name}} Details :
-----------------
NAME            : {{.Name}}
ACCESS MODE     : {{.AccessMode}}
CSI DRIVER      : {{.CSIDriver}}
STORAGE CLASS   : {{.StorageClass}}
VOLUME PHASE    : {{.VolumePhase }}
VERSION         : {{.Version}}
CSPC            : {{.CSPC}}
SIZE            : {{.Size}}
STATUS          : {{.Status}}
REPLICA COUNT	: {{.ReplicaCount}}

`

	portalTemplate = `
Portal Details :
------------------
IQN              :  {{.IQN}}
VOLUME NAME      :  {{.VolumeName}}
TARGET NODE NAME :  {{.TargetNodeName}}
PORTAL           :  {{.Portal}}
TARGET IP        :  {{.TargetIP}}

`
)

// NewCmdDescribeVolume displays OpenEBS Volume information.
func NewCmdDescribeVolume() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "volume",
		Aliases: []string{"volumes", "vol", "v", "vols"},
		Short:   "Displays Openebs volume information",
		Long:    volumeInfoCommandHelpText,
		Example: `kubectl openebs describe volume [vol]`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Get this from flags, pflag, etc
			openebsNs, _ := cmd.Flags().GetString("openebs-namespace")
			util.CheckErr(RunVolumeInfo(cmd, args, openebsNs), util.Fatal)
		},
	}
	return cmd
}

// RunVolumeInfo runs info command and make call to DisplayVolumeInfo to display the results
func RunVolumeInfo(cmd *cobra.Command, vols []string, openebsNs string) error {
	// the stuff automatically coming from kubectl command execution
	clientset, err := client.NewK8sClient(openebsNs)
	util.CheckErr(err, util.Fatal)

	if openebsNs == "" {
		nsFromCli, err := clientset.GetOpenEBSNamespace(util.CstorCasType)
		if err != nil {
			//return errors.Wrap(err, "Error determining the openebs namespace, please specify using \"--openebs-namespace\" flag")
			return errors.New("no cstor volumes found in the cluster")
		}
		clientset.Ns = nsFromCli
	}

	// TODO: Print all volume info present in args or print all volume info if no args given
	if len(vols) == 0 {
		return errors.New("Please give at least one volume to describe")
	}
	for _, volName := range vols {
		// Fetch all details of a volume is called to get the volume controller's
		// info such as controller's IP, status, iqn, replica IPs etc.
		//1. cStor volume info
		volumeInfo, err := clientset.GetcStorVolume(volName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get CStorVolume %s\n", volName)
			continue
		}
		//2. Persistent Volume info
		pvInfo, err := clientset.GetPV(volName)
		if err != nil {
			fmt.Printf("failed to get PV %s\n", volName)
			continue
		}
		//3. cStor Volume Config
		cvcInfo, err := clientset.GetCVC(volName)
		if err != nil {
			fmt.Printf("failed to get cStor Volume config for %s", volName)
			continue
		}

		//4. Get Node for Target Pod from the openebs-ns
		node, err := clientset.GetCStorVolumeAttachment(volName)
		var nodeName string
		if err != nil {
			nodeName = util.NotAttached
			fmt.Fprintf(os.Stderr, "failed to get CStorVolumeAttachments for %s\n", volName)
		} else {
			nodeName = node.Spec.Volume.OwnerNodeID
		}

		//5. cStor Volume Replicas
		cvrInfo, err := clientset.GetCVR(volName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get cStor Volume Replicas for %s\n", volName)
		}
		cSPCLabel := cstortypes.CStorPoolClusterLabelKey
		volume := util.VolumeInfo{
			AccessMode:              util.AccessModeToString(pvInfo.Spec.AccessModes),
			Capacity:                volumeInfo.Status.Capacity.String(),
			CSPC:                    cvcInfo.Labels[cSPCLabel],
			CSIDriver:               pvInfo.Spec.CSI.Driver,
			CSIVolumeAttachmentName: pvInfo.Spec.CSI.VolumeHandle,
			Name:                    volumeInfo.Name,
			Namespace:               volumeInfo.Namespace,
			PVC:                     pvInfo.Spec.ClaimRef.Name,
			ReplicaCount:            volumeInfo.Spec.ReplicationFactor,
			VolumePhase:             pvInfo.Status.Phase,
			StorageClass:            pvInfo.Spec.StorageClassName,
			Version:                 util.CheckVersion(volumeInfo.VersionDetails),
			Size:                    util.ConvertToIBytes(volumeInfo.Status.Capacity.String()),
			Status:                  volumeInfo.Status.Phase,
		}

		// Print the output for the portal status info
		err = util.PrintByTemplate("volume", volInfoTemplate, volume)
		if err != nil {
			return err
		}

		portalInfo := util.PortalInfo{
			IQN:            volumeInfo.Spec.Iqn,
			VolumeName:     volumeInfo.Name,
			Portal:         volumeInfo.Spec.TargetPortal,
			TargetIP:       volumeInfo.Spec.TargetIP,
			TargetNodeName: nodeName,
		}

		// Print the output for the portal status info
		err = util.PrintByTemplate("PortalInfo", portalTemplate, portalInfo)
		if err != nil {
			return err
		}

		replicaCount := volumeInfo.Spec.ReplicationFactor
		// This case will occur only if user has manually specified zero replica.
		// or if none of the replicas are healthy & running
		if replicaCount == 0 || len(volumeInfo.Status.ReplicaStatuses) == 0 {
			fmt.Printf("None of the replicas are running\n")
			//please check the volume pod's status by running [kubectl describe pvc -l=openebs/replica --all-namespaces]\Oor try again later.")
			return nil
		}

		// Print replica details
		if cvrInfo != nil && len(cvrInfo.Items) > 0 {
			fmt.Printf("\nReplica Details :\n-----------------\n")
			var rows []metav1.TableRow
			for _, cvr := range cvrInfo.Items {
				rows = append(rows, metav1.TableRow{Cells: []interface{}{
					cvr.Name,
					util.ConvertToIBytes(cvr.Status.Capacity.Total),
					util.ConvertToIBytes(cvr.Status.Capacity.Used),
					cvr.Status.Phase,
					util.Duration(time.Since(cvr.ObjectMeta.CreationTimestamp.Time))}})
			}
			util.TablePrinter(util.CstorReplicaColumnDefinations, rows, printers.PrintOptions{Wide: true})
		}

		cStorBackupList, err := clientset.GetCstorVolumeBackups(volName)
		if cStorBackupList != nil {
			fmt.Printf("\nCstor Backup Details :\n" + "---------------------\n")
			var rows []metav1.TableRow
			for _, item := range cStorBackupList.Items {
				rows = append(rows, metav1.TableRow{Cells: []interface{}{
					item.ObjectMeta.Name,
					item.Spec.BackupName,
					item.Spec.VolumeName,
					item.Spec.BackupDest,
					item.Spec.SnapName,
					item.Status,
				}})
			}
			util.TablePrinter(util.CstorBackupColumnDefinations, rows, printers.PrintOptions{Wide: true})
		}

		cstorCompletedBackupList, err := clientset.GetCstorVolumeCompletedBackups(volName)
		if cstorCompletedBackupList != nil {
			fmt.Printf("\nCstor Completed Backup Details :" + "\n-------------------------------\n")
			var rows []metav1.TableRow
			for _, item := range cstorCompletedBackupList.Items {
				rows = append(rows, metav1.TableRow{Cells: []interface{}{
					item.Name,
					item.Spec.BackupName,
					item.Spec.VolumeName,
					item.Spec.LastSnapName,
				}})
			}
			util.TablePrinter(util.CstorCompletedBackupColumnDefinations, rows, printers.PrintOptions{Wide: true})
		}

		cStorRestoreList, err := clientset.GetCstorVolumeRestores(volName)
		if cStorRestoreList != nil {
			fmt.Printf("\nCstor Restores Details :" + "\n-----------------------\n")
			var rows []metav1.TableRow
			for _, item := range cStorRestoreList.Items {
				rows = append(rows, metav1.TableRow{Cells: []interface{}{
					item.ObjectMeta.Name,
					item.Spec.RestoreName,
					item.Spec.VolumeName,
					item.Spec.RestoreSrc,
					item.Spec.StorageClass,
					item.Spec.Size.String(),
					item.Status,
				}})
			}
			util.TablePrinter(util.CstorRestoreColumnDefinations, rows, printers.PrintOptions{Wide: true})
		}
		fmt.Println()
	}
	return nil
}
