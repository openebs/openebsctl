/*
Copyright 2020 The OpenEBS Authors

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

package volume

import (
	"fmt"
	"html/template"
	"os"

	cstortypes "github.com/openebs/api/pkg/apis/types"

	"github.com/openebs/openebsctl/client"
	"github.com/openebs/openebsctl/kubectl-openebs/cli/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"k8s.io/klog"
)

var (
	volumeInfoCommandHelpText = `
This command fetches information and status of the various
aspects of a cStor Volume such as ISCSI, Controller, and Replica.

Usage: openebs cstor volume describe --volname <vol>
`

	volName string
)

const (
	volInfoTemplate = `
Volume Details :
----------------
Name            : {{.Name}}
Access Mode     : {{.AccessMode}}
CSI Driver      : {{.CSIDriver}}
Storage Class   : {{.StorageClass}}
Volume Phase    : {{.VolumePhase }}
Version         : {{.Version}}
CSPC            : {{.CSPC}}
Size            : {{.Size}}
Status          : {{.Status}}
ReplicaCount	: {{.ReplicaCount}}
`

	portalTemplate = `
Portal Details :
----------------
IQN             :  {{.IQN}}
Volume          :  {{.VolumeName}}
TargetNodeName  :  {{.TargetNodeName}}
Portal          :  {{.Portal}}
TargetIP        :  {{.TargetIP}}

`
)

// NewCmdVolumeInfo displays OpenEBS Volume information.
func NewCmdVolumeInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "describe",
		Short:   "Displays Openebs Volume information",
		Long:    volumeInfoCommandHelpText,
		Example: `openebs cStor volume describe --volname <vol>`,
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(RunVolumeInfo(cmd), util.Fatal)
		},
	}
	cmd.Flags().StringVarP(&volName, "volname", "", volName,
		"unique volume name.")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "openebs",
		"namespace name, required if the OpenEBS is not in the openebs namespace")

	return cmd
}

// RunVolumeInfo runs info command and make call to DisplayVolumeInfo to display the results
func RunVolumeInfo(cmd *cobra.Command) error {

	clientset, err := client.NewK8sClient(namespace)
	if err != nil {
		return errors.Wrap(err, "failed to execute volume info command")
	}

	// Fetch all details of a volume is called to get the volume controller's
	// info such as controller's IP, status, iqn, replica IPs etc.
	//1. cStor volume info
	volumeInfo, err := clientset.GetcStorVolume(volName)
	if err != nil {
		return errors.Wrap(err, "failed to execute volume info command, getting cStor volumes")
	}
	//2. Persistent Volume info
	pvInfo, err := clientset.GetPV(volName)
	if err != nil {
		return errors.Wrap(err, "failed to execute volume info command, getting persistant volumes")
	}

	//3. cStor Volume Config
	cvcInfo, err := clientset.GetCVC(volName)
	if err != nil {
		return errors.Wrap(err, "failed to execute volume info command, getting cStor Volume config")
	}

	//4. Get Node Name for Target Pod
	NodeName, err := clientset.NodeForVolume(volName)
	if err != nil {
		klog.Errorf("error executeing volume info command, getting Node for Volume %s:{%s}", volName, err)
	}

	//5. cStor Volume Replicas
	cvrInfo, err := clientset.GetCVR(volName)
	if err != nil {
		return errors.Wrap(err, "failed to execute volume info command, getting cStor Volume Replicas")
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
		Size:                    volumeInfo.Status.Capacity.String(),
		Status:                  volumeInfo.Status.Phase,
	}

	// Print the output for the portal status info
	tmpl, err := template.New("volume").Parse(volInfoTemplate)
	if err != nil {
		return errors.Wrap(err, "error displaying output for volume info")
	}
	err = tmpl.Execute(os.Stdout, volume)
	if err != nil {
		return errors.Wrap(err, "error displaying volume details")

	}

	portalInfo := util.PortalInfo{
		IQN:            volumeInfo.Spec.Iqn,
		VolumeName:     volumeInfo.Name,
		Portal:         volumeInfo.Spec.TargetPortal,
		TargetIP:       volumeInfo.Spec.TargetIP,
		TargetNodeName: NodeName,
	}

	// Print the output for the portal status info
	tmpl, err = template.New("PortalInfo").Parse(portalTemplate)
	if err != nil {
		return errors.Wrap(err, "error creating output for portal info")
	}
	err = tmpl.Execute(os.Stdout, portalInfo)
	if err != nil {
		fmt.Println(err, "error displaying target portal detail")
		return nil
	}

	replicaCount := volumeInfo.Spec.ReplicationFactor
	// This case will occur only if user has manually specified zero replica.
	// or if none of the replicas are healthy & running
	if replicaCount == 0 || len(volumeInfo.Status.ReplicaStatuses) == 0 {
		fmt.Println("None of the replicas are running")
		//please check the volume pod's status by running [kubectl describe pvc -l=openebs/replica --all-namespaces]\Oor try again later.")
		return nil
	}

	// Print replica details
	fmt.Printf("Replica Details :\n----------------\n")
	out := make([]string, len(cvrInfo.Items)+2)
	out[0] = "Name|Pool Instance|Status"
	out[1] = "----|-------------|------"
	for i, cvr := range cvrInfo.Items {
		out[i+2] = fmt.Sprintf("%s|%s|%s",
			cvr.ObjectMeta.Name,
			cvr.Labels[cstortypes.CStorPoolInstanceNameLabelKey],
			cvr.Status.Phase,
		)
	}

	fmt.Println(util.FormatList(out))
	return nil
}
