package volume

import (
	"flag"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vaniisgh/mayactl/client"
	"github.com/vaniisgh/mayactl/kubectl-mayactl/cli/util"
)

var (
	volumesListCommandHelpText = `
This command displays status of available zfs Volumes.
If no volume ID is given, a list of all known volumes will be displayed.

Usage: kubectl mayactl volume list [options]
	`
	namespace string
)

const (
	volumeStatusOK = "Running"
)

// NewCmdVolumesList displays status of OpenEBS Volume(s)
func NewCmdVolumesList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Displays status information about Volume(s)",
		Long:  volumesListCommandHelpText,
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(RunVolumesList(cmd), util.Fatal)
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "openebs",
		"namespace name, required if volume is not in the `default` namespace")

	flag.CommandLine.Parse([]string{})

	return cmd
}

//RunVolumesList fetchs the volumes from maya-apiserver
func RunVolumesList(cmd *cobra.Command) error {

	client, err := client.NewK8sClient(namespace)
	util.CheckErr(err, util.Fatal)

	cvols := client.GetcStorVolumes()
	pvols := client.GetcStorPVCs("")

	// tally status of cvols to pvols
	//give output according to volume status
	out := make([]string, len(cvols.Items)+2)
	out[0] = "Node|Namespace|Name|csiVolumeAttachmentName|Status|Type|Version|Capacity|StorageClass|Attached|Access Mode"
	out[1] = "----|---------|----|-----------------------|------|----|-------|--------|------------|--------|-----------"
	for i, item := range cvols.Items {
		out[i+2] = fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|%s",
			pvols[item.ObjectMeta.Name].Node,
			item.ObjectMeta.Namespace,
			item.ObjectMeta.Name,
			pvols[item.ObjectMeta.Name].CSIVolumeAttachmentName,
			item.Status.Phase,
			item.TypeMeta.Kind,
			item.VersionDetails.Status.Current,
			item.Status.Capacity.String(),
			pvols[item.ObjectMeta.Name].StorageClass,
			pvols[item.ObjectMeta.Name].AttachementStatus,
			pvols[item.ObjectMeta.Name].AccessMode,
		)
	}
	if len(out) == 2 {
		fmt.Println("No Volumes are running")
		return nil
	}

	fmt.Println(util.FormatList(out))

	return nil
}
