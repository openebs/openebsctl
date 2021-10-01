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

package upgrade

import (
	"fmt"
	"os"

	"github.com/openebs/openebsctl/pkg/upgrade"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/spf13/cobra"
)

const (
	upgradeCmdHelp = `Upgrade OpenEBS Data Plane Components
  
	Usage: 
    kubectl openebs upgrade volume [flags]

  Flags:
  -h, --help                   help for openebs upgrade command
  -f, --file                   provide menifest file containing job upgrade information
      --openebs-namespace      upgrade by a fixed openEBS namespace
      --cas-type               [jiva | cStor | LocalPv] specify the cas-type to upgrade
      --to-version             the desired version for upgradation
      --image-prefix           if required the image prefix of the volume deployments can be
                               changed using the flag, defaults to whatever was present on old
                               deployments.
      --image-tag              if required the image tags for volume deployments can be changed
                               to a custom image tag using the flag, 
                               defaults to the --to-version mentioned above.
	`
)

// NewCmdClusterInfo to upgrade volumes and interfaces
func NewCmdVolumeUpgrade(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "upgrade",
		Short:   "Upgrade Volumes, storage engines, and interfaces in openebs application",
		Aliases: []string{"update"},
		Run: func(cmd *cobra.Command, args []string) {
			openebsNs, _ := cmd.Flags().GetString("openebs-namespace")
			if !util.IsValidCasType(upgrade.CasType) {
				fmt.Fprintf(os.Stderr, "cas-type %s not supported\n", upgrade.CasType)
			} else if upgrade.CasType == util.JivaCasType {
				upgrade.InstantiateJivaUpgrade(openebsNs)
			} else {
				fmt.Println("Only Jiva upgrades are available at this point")
				fmt.Println("To upgrade other cas-types follow: https://github.com/openebs/upgrade#upgrading-openebs-reources")
			}
		},
	}
	cmd.SetUsageTemplate(upgradeCmdHelp)
	cmd.PersistentFlags().StringVarP(&upgrade.CasType, "cas-type", "", "", "the cas-type filter option for fetching resources")
	cmd.PersistentFlags().StringVarP(&upgrade.ToVersion, "to-version", "", "", "the version to which the resources need to be upgraded")
	cmd.PersistentFlags().StringVarP(&upgrade.File, "file", "f", "", "provide path/url to the menifest file")
	cmd.PersistentFlags().StringVarP(&upgrade.ImagePrefix, "image-prefix", "", "", "provide image prefix for the volume deployments")
	cmd.PersistentFlags().StringVarP(&upgrade.ImageTag, "image-tag", "", "", "provide custom image tag for the volume deployments")
	return cmd
}
