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
	// "github.com/openebs/openebsctl/pkg/upgrade"
)

const (
	upgradeCmdHelp = `Upgrade OpenEBS Data Plane Components
  
	Usage: 
    kubectl openebs upgrade volume [flags]

  Flags:
  -h, --help                   help for openebs upgrade command
	-f,	--file                   provide menifest file containing job upgrade information
	    --openebs-namespace      upgrade by a fixed openEBS namespace
			--cas-type               [jiva | cStor | LocalPv] specify the cas-type
			                         to upgrade
		  --to-version             the desired version for upgradation
	`
)

// NewCmdClusterInfo to upgrade volumes and interfaces
func NewCmdVolumeUpgrade(rootCmd *cobra.Command) *cobra.Command {
	var casType, toVersion, file string
	cmd := &cobra.Command{
		Use:     "upgrade",
		Short:   "Upgrade Volumes, storage engines, and interfaces in openebs application",
		Aliases: []string{"update"},
		Run: func(cmd *cobra.Command, args []string) {
			openebsNs, _ := cmd.Flags().GetString("openebs-namespace")
			cType, _ := cmd.Flags().GetString("cas-type")
			toVersion, _ := cmd.Flags().GetString("to-version")
			menifestFile, _ := cmd.Flags().GetString("file")
			if !util.IsValidCasType(cType) {
				fmt.Fprintf(os.Stderr, "cas-type %s not supported\n", cType)
			} else if cType == util.JivaCasType {
				upgrade.InstantiateJivaUpgrade(openebsNs, toVersion, menifestFile)
			} else {
				fmt.Println("Only Jiva upgrades are available at this point")
				fmt.Println("To upgrade other cas-types follow: https://github.com/openebs/upgrade#upgrading-openebs-reources")
			}
		},
	}
	cmd.SetUsageTemplate(upgradeCmdHelp)
	cmd.PersistentFlags().StringVarP(&casType, "cas-type", "", "", "the cas-type filter option for fetching resources")
	cmd.PersistentFlags().StringVarP(&toVersion, "to-version", "", "", "the version to which the resources need to be upgraded")
	cmd.PersistentFlags().StringVarP(&file, "file", "f", "", "provide path/url to the menifest file")
	return cmd
}
