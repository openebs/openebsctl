package volume

import (
	"fmt"
	"html/template"
	"os"

	//"strconv"

	cstorv1 "github.com/openebs/api/pkg/apis/cstor/v1"

	"github.com/spf13/cobra"
	"github.com/vaniisgh/mayactl/client"
	"github.com/vaniisgh/mayactl/kubectl-mayactl/cli/util"
)

var (
	volumeInfoCommandHelpText = `
This command fetches information and status of the various
aspects of a cStor Volume such as ISCSI, Controller, and Replica.

Usage: mayactl cstor volume describe --volname <vol>
`

	cStorVolumeSpec   *cstorv1.CStorVolumeSpec
	cStorVolumeStatus *cstorv1.CStorVolumeStatus
	versionDetails    *cstorv1.VersionDetails
	metaDataDetails   *metaDataInfo
	volName           string
)

type metaDataInfo struct {
	test string
}

const (
	volInfoTemplate = `
Replica Details :
-----------------
{{ printf "NAME\t ACCESSMODE\t STATUS\t IP\t NODE" }}
{{ printf "-----\t -----------\t -------\t ---\t -----" }} {{range $key, $value := .}}
{{ printf "%s\t" $value.Name }} {{ printf "%s\t" $value.AccessMode }} {{ printf "%s\t" $value.Status }} {{ printf "%s\t" $value.IP }} {{ $value.NodeName }} {{end}}
`

	cstorReplicaTemplate = `
Replica Details :
-----------------
{{ printf "%s\t" "NAME"}} {{ printf "%s\t" "STATUS"}} {{ printf "%s\t" "POOL NAME"}} {{ printf "%s\t" "NODE"}}
{{ printf "----\t ------\t ---------\t -----" }} {{range $key, $value := .}}
{{ printf "%s\t" $value.Name }} {{ printf "%s\t" $value.Status }} {{ printf "%s\t" $value.PoolName }} {{ $value.NodeName }} {{end}}
`

	portalTemplate = `
Portal Details :
----------------
IQN               :   {{.IQN}}
Volume            :   {{.VolumeName}}
Portal            :   {{.Portal}}
Size              :   {{.Size}}
Controller Status :   {{.Status}}
Replica Count     :   {{.ReplicaCount}}
Replica Status     :   {{.ReplicaStatus}}
`

//Controller Node   :   {{.ControllerNode}}

)

// NewCmdVolumeInfo displays OpenEBS Volume information.
func NewCmdVolumeInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "describe",
		Short:   "Displays Openebs Volume information",
		Long:    volumeInfoCommandHelpText,
		Example: `mayactl cStor volume describe --volname <vol>`,
		Run: func(cmd *cobra.Command, args []string) {
			//util.CheckErr(errors.Validate(cmd, false, false, true), util.Fatal)
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

	client, err := client.NewK8sClient(namespace)

	// Fetch all details of a volume is called to get the volume controller's
	// info such as controller's IP, status, iqn, replica IPs etc.
	volumeInfo := client.GetcStorVolume(volName, namespace)

	portalInfo := util.PortalInfo{
		volumeInfo.Spec.Iqn,
		volumeInfo.Name,
		volumeInfo.Spec.TargetPortal,
		volumeInfo.Status.Capacity.String(),
		volumeInfo.Status.Conditions,
		volumeInfo.Spec.ReplicationFactor,
		volumeInfo.Status.ReplicaStatuses,
	}

	// Print the output for the portal status info
	tmpl, err := template.New("VolumeInfo").Parse(portalTemplate)
	if err != nil {
		fmt.Println("Error displaying output, found error :", err)
		return nil
	}
	err = tmpl.Execute(os.Stdout, portalInfo)
	if err != nil {
		fmt.Println("Error displaying volume details, found error :", err)
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
	// Splitting strings with delimiter ','
	//replicaStatusStrings := strings.Split(volumeInfo.Status.Message, ",")
	//addressIPStrings := strings.Split(volumeInfo.Spec.TargetIP, ",")

	return nil
}
