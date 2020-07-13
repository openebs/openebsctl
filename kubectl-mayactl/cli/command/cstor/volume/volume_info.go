package volume

/*
import (
	"fmt"
	"html/template"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vaniisgh/mayactl/client"
	"github.com/vaniisgh/mayactl/kubectl-mayactl/cli/util"

	cstorv1 "github.com/openebs/api/pkg/apis/cstor/v1"
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
Controller Node   :   {{.ControllerNode}}
Replica Count     :   {{.ReplicaCount}}
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
			//util.CheckErr(Validate(cmd, false, false, true), util.Fatal)
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
	util.CheckErr(err, util.Fatal)

	// Fetch all details of a volume is called to get the volume controller's
	// info such as controller's IP, status, iqn, replica IPs etc.
	volumeInfo := client.GetcStorVolume(volName, namespace)

	portalInfo = PortalInfo{
		v.GetIQN(),
		v.GetVolumeName(),
		v.GetTargetPortal(),
		v.GetVolumeSize(),
		v.GetControllerStatus(),
		v.GetReplicaCount(),
		v.GetControllerNode(),
	}

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

	if v.GetCASType() == string(JivaStorageEngine) {
		replicaCount, _ = strconv.Atoi(v.GetReplicaCount())
		// This case will occur only if user has manually specified zero replica.
		if replicaCount == 0 || len(v.GetReplicaStatus()) == 0 {
			fmt.Println("None of the replicas are running, please check the volume pod's status by running [kubectl describe pod -l=openebs/replica --all-namespaces] or try again later.")
			return nil
		}
		// Splitting strings with delimiter ','
		replicaStatusStrings := strings.Split(v.GetReplicaStatus(), ",")
		addressIPStrings := strings.Split(v.GetReplicaIP(), ",")

		// making a map of replica ip and their respective status,index and mode
		replicaIPStatus := make(map[string]*Value)

		// Creating a map of address and mode. The IP is chosen as key so that the status of that corresponding replica can be merged in linear time complexity
		for index, IP := range addressIPStrings {
			if strings.Contains(IP, "nil") {
				// appending address with index to avoid same key conflict as the IP is returned as `nil` in case of error
				replicaIPStatus[IP+string(index)] = &Value{index: index, status: replicaStatusStrings[index], mode: "NA"}
			} else {
				replicaIPStatus[IP] = &Value{index: index, status: replicaStatusStrings[index], mode: "NA"}
			}
		}

		// We get the info of the running replicas from the collection.data.
		// We are appending modes if available in collection.data to replicaIPStatus
		replicaInfo := make(map[int]*ReplicaInfo)

		for key := range collection.Data {
			address = append(address, strings.TrimSuffix(strings.TrimPrefix(collection.Data[key].Address, "tcp://"), v1.ReplicaPort))
			mode = append(mode, collection.Data[key].Mode)
			if _, ok := replicaIPStatus[address[key]]; ok {
				replicaIPStatus[address[key]].mode = mode[key]
			}
		}

		for IP, replicaStatus := range replicaIPStatus {
			// checking if the first three letters is nil or not if it is nil then the ip is not available
			if strings.Contains(IP, "nil") {
				replicaInfo[replicaStatus.index] = &ReplicaInfo{"NA", replicaStatus.mode, replicaStatus.status, "NA", "NA"}
			} else {
				replicaInfo[replicaStatus.index] = &ReplicaInfo{IP, replicaStatus.mode, replicaStatus.status, "NA", "NA"}
			}
		}
		// updating the replica info to replica structure
		err = updateReplicasInfo(replicaInfo)
		if err != nil {
			fmt.Println("Error in getting specific information from K8s. Please try again.")
		}

		return mapiserver.Print(jivaReplicaTemplate, replicaInfo)
	} else if v.GetCASType() == string(CstorStorageEngine) {

		// Converting replica count character to int
		replicaCount, err = strconv.Atoi(v.GetReplicaCount())
		if err != nil {
			fmt.Println("Invalid replica count")
			return nil
		}

		// Spitting the replica status
		replicaStatus := strings.Split(v.GetControllerStatus(), ",")
		poolName := strings.Split(v.GetStoragePool(), ",")
		cvrName := strings.Split(v.GetCVRName(), ",")
		nodeName := strings.Split(v.GetNodeName(), ",")

		// Confirming replica status, poolname , cvrName, nodeName are equal to replica count
		//if replicaCount != len(replicaStatus) || replicaCount != len(poolName) || replicaCount != len(cvrName) || replicaCount != len(nodeName) {
		if replicaCount != len(poolName) || replicaCount != len(cvrName) || replicaCount != len(nodeName) {
			fmt.Println("Invalid response received from maya-api service")
			return nil
		}

		replicaInfo := []cstorReplicaInfo{}

		// Iterating over the values replica values and appending to the structure
		for i := 0; i < replicaCount; i++ {
			replicaInfo = append(replicaInfo, cstorReplicaInfo{
				Name:       cvrName[i],
				PoolName:   poolName[i],
				AccessMode: "N/A",
				Status:     strings.Title(replicaStatus[i]),
				NodeName:   nodeName[i],
				IP:         "N/A",
			})
		}

		return mapiserver.Print(cstorReplicaTemplate, replicaInfo)
	} else {
		fmt.Println("Unsupported Volume Type")
	}
	return nil
}
*/
