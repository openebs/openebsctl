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

package get

import (
	"fmt"

	v1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
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
)

// NewCmdGetVolume displays status of OpenEBS Volume(s)
func NewCmdGetVolume() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "volume",
		Aliases: []string{"vol", "v", "volumes"},
		Short:   "Displays status information about Volume(s)",
		Long:    volumesListCommandHelpText,
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(RunVolumesList(cmd, args), util.Fatal)
		},
	}
	return cmd
}

// RunVolumesList lists the volumes
func RunVolumesList(cmd *cobra.Command, vols []string) error {
	client, err := client.NewK8sClient("")
	util.CheckErr(err, util.Fatal)
	var cvols *v1.CStorVolumeList
	if len(vols) == 0 {
		cvols, err = client.GetcStorVolumes()
	} else {
		cvols, err = client.GetcStorVolumesByNames(vols)
	}
	if err != nil {
		return errors.Wrap(err, "error listing volumes")
	}
	pvols, err := client.GetCStorVolumeInfoMap("")
	if err != nil {
		return errors.Wrap(err, "failed to execute volume info command")
	}
	// tally status of cvols to pvols
	// give output according to volume status
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
			pvols[item.ObjectMeta.Name].Node)
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
