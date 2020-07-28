package volume

import (
	"flag"
	"fmt"

	"github.com/openebs/openebsctl/client"
	"github.com/openebs/openebsctl/kubectl-openebs/cli/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	volumesListCommandHelpText = `
This command displays status of available zfs Volumes.
If no volume ID is given, a list of all known volumes will be displayed.

Usage: kubectl openebs cStor volume list [options]
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

	cvols, err := client.GetcStorVolumes()
	if err != nil {
		return errors.Wrap(err, "error listing volumes")
	}
	pvols, err := client.GetCStorVolumeInfoMap("")
	if err != nil {
		return errors.Wrap(err, "failed to execute volume info command")
	}

	// tally status of cvols to pvols
	//give output according to volume status
	out := make([]string, len(cvols.Items)+2)
	out[0] = "Namespace|Name|Status|Version|Capacity|StorageClass|Attached|Access Mode|Attached Node"
	out[1] = "---------|----|------|-------|--------|------------|--------|-----------|-------------"
	for i, item := range cvols.Items {
		pvols[item.ObjectMeta.Name] = util.CheckForVol(item.ObjectMeta.Name, pvols)
		out[i+2] = fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s",
			item.ObjectMeta.Namespace,
			item.ObjectMeta.Name,
			item.Status.Phase,
			item.VersionDetails.Status.Current,
			item.Status.Capacity.String(),
			pvols[item.ObjectMeta.Name].StorageClass,
			pvols[item.ObjectMeta.Name].AttachementStatus,
			pvols[item.ObjectMeta.Name].AccessMode,
			pvols[item.ObjectMeta.Name].Node,
		)

		//TODO: find a fix
		//pvols[item.ObjectMeta.Name].CSIVolumeAttachmentName field removed for readability
	}
	if len(out) == 2 {
		fmt.Println("No Volumes are running")
		return nil
	}

	fmt.Println(util.FormatList(out))

	return nil
}
