package volume

import (
	"fmt"
	"html/template"
	"os"

	apiTypes "github.com/openebs/api/pkg/apis/types"

	"github.com/spf13/cobra"
	"github.com/vaniisgh/mayactl/client"
	"github.com/vaniisgh/mayactl/kubectl-mayactl/cli/util"

	"k8s.io/klog"
)

var (
	volumeInfoCommandHelpText = `
This command fetches information and status of the various
aspects of a cStor Volume such as ISCSI, Controller, and Replica.

Usage: mayactl cstor volume describe --volname <vol>
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
		Example: `mayactl cStor volume describe --volname <vol>`,
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(RunVolumeInfo(cmd), util.Fatal)
		},
	}
	cmd.Flags().StringVarP(&volName, "volname", "", volName,
		"unique volume name.")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "openebs",
		"namespace name, required if volume is not in the `default` namespace")

	return cmd
}

// RunVolumeInfo runs info command and make call to DisplayVolumeInfo to display the results
func RunVolumeInfo(cmd *cobra.Command) error {

	clientset, err := client.NewK8sClient(namespace)

	// Fetch all details of a volume is called to get the volume controller's
	// info such as controller's IP, status, iqn, replica IPs etc.
	//1. cStor volume info
	volumeInfo := clientset.GetcStorVolume(volName)
	//2. Persistent Volume info
	pvInfo := clientset.GetPV(volName)
	//3. cStor Volume Config
	cvcInfo := clientset.GetCVC(volName)
	//4. Get Node Name for Target Pod
	NodeName := clientset.NodeForVolume(volName)
	//5. cStor Volume Replicas
	cvrInfo := clientset.GetCVR(volName)

	cSPCLabel := apiTypes.CStorPoolClusterLabelKey

	volume := util.VolumeInfo{
		util.AccessModeToString(pvInfo.Spec.AccessModes),
		volumeInfo.Status.Capacity.String(),
		cvcInfo.Labels[cSPCLabel],
		pvInfo.Spec.CSI.Driver,
		pvInfo.Spec.CSI.VolumeHandle,
		volumeInfo.Name,
		volumeInfo.Namespace,
		pvInfo.Spec.ClaimRef.Name,
		volumeInfo.Spec.ReplicationFactor,
		pvInfo.Status.Phase,
		pvInfo.Spec.StorageClassName,
		util.CheckVersion(volumeInfo.VersionDetails),
		volumeInfo.Status.Capacity.String(),
		volumeInfo.Status.Phase,
	}

	// Print the output for the portal status info
	tmpl, err := template.New("volume").Parse(volInfoTemplate)
	if err != nil {
		fmt.Println("Error displaying output, found error :", err)
		return nil
	}
	err = tmpl.Execute(os.Stdout, volume)
	if err != nil {
		fmt.Println("Error displaying volume details, found error :", err)
		return nil
	}

	portalInfo := util.PortalInfo{
		volumeInfo.Spec.Iqn,
		volumeInfo.Name,
		volumeInfo.Spec.TargetPortal,
		volumeInfo.Spec.TargetIP,
		NodeName,
	}

	// Print the output for the portal status info
	tmpl, err = template.New("PortalInfo").Parse(portalTemplate)
	if err != nil {
		fmt.Println("Error displaying output, found error :", err)
		return nil
	}
	err = tmpl.Execute(os.Stdout, portalInfo)
	if err != nil {
		fmt.Println("Error displaying protal details, found error :", err)
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

	//map of replica ID to replica
	knownReplicas := volumeInfo.Status.ReplicaDetails.KnownReplicas

	if len(knownReplicas) != replicaCount {
		klog.Errorf("WARNING: Desired replica count %s while known replicas are %s", replicaCount, knownReplicas)
	}

	cSPILabel := apiTypes.CStorPoolInstanceNameLabelKey

	fmt.Printf("Replica Details :\n----------------\n")
	out := make([]string, len(cvrInfo)+2)
	out[0] = "Name|Pool Instance|Status"
	out[1] = "----|-------------|------"
	for i, cvr := range cvrInfo {
		out[i+2] = fmt.Sprintf("%s|%s|%s",
			cvr.ObjectMeta.Name,
			cvr.Labels[cSPILabel],
			cvr.Status.Phase,
		)
	}

	fmt.Println(util.FormatList(out))
	return nil
}
