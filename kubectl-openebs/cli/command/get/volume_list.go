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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"

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

Usage: kubectl openebs get volume [options]
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
			openebsNs, _ := cmd.Flags().GetString("openebs-namespace")
			util.CheckErr(RunVolumesList(cmd, openebsNs, args), util.Fatal)
		},
	}
	return cmd
}

// RunVolumesList lists the volumes
func RunVolumesList(cmd *cobra.Command, openebsNs string, vols []string) error {
	k8sClient, err := client.NewK8sClient("")
	util.CheckErr(err, util.Fatal)
	if openebsNs == "" {
		nsFromCli, err := k8sClient.GetOpenEBSNamespace(util.CstorCasType)
		if err != nil {
			return errors.Wrap(err, "Error determining the openebs namespace, please specify using \"--openebs-namespace\" flag")
		}
		k8sClient.Ns = nsFromCli
	}
	var cvols *v1.CStorVolumeList
	if len(vols) == 0 {
		cvols, err = k8sClient.GetcStorVolumes()
	} else {
		cvols, err = k8sClient.GetcStorVolumesByNames(vols)
	}
	if err != nil {
		return errors.Wrap(err, "error listing volumes")
	}
	pvols, err := k8sClient.GetCStorVolumeInfoMap("")
	if err != nil {
		return errors.Wrap(err, "failed to execute volume info command")
	}
	// tally status of cvols to pvols
	// give output according to volume status
	var rows []metav1.TableRow
	for _, item := range cvols.Items {
		rows = append(rows, metav1.TableRow{Cells: []interface{}{
			item.ObjectMeta.Namespace,
			item.ObjectMeta.Name,
			item.Status.Phase,
			item.VersionDetails.Status.Current,
			util.ConvertToIBytes(item.Status.Capacity.String()),
			pvols[item.ObjectMeta.Name].StorageClass,
			pvols[item.ObjectMeta.Name].AttachementStatus,
			pvols[item.ObjectMeta.Name].AccessMode,
			pvols[item.ObjectMeta.Name].Node}})
		//TODO: find a fix
		//pvols[item.ObjectMeta.Name].CSIVolumeAttachmentName field removed for readability
	}
	util.TablePrinter(util.CstorVolumeListColumnDefinations, rows, printers.PrintOptions{Wide: true})
	return nil
}
