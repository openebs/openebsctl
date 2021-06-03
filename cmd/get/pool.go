package get

import (
	"github.com/openebs/openebsctl/pkg/cstor/get"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/spf13/cobra"
)

var (
	poolListCommandHelpText = `
This command lists of all known pools in the Cluster.

Usage:
$ kubectl openebs get pool [options]
`
)

// NewCmdGetPool displays status of OpenEBS Pool(s)
func NewCmdGetPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pool",
		Aliases: []string{"pools", "p"},
		Short:   "Displays status information about Pool(s)",
		Long:    poolListCommandHelpText,
		Run: func(cmd *cobra.Command, args []string) {
			ns, err := cmd.Flags().GetString("namespace")
			if err != nil {
				ns = "openebs"
			}
			// TODO: De-couple CLI code, logic code, API code
			util.CheckErr(get.RunPoolsList(cmd, ns), util.Fatal)
		},
	}
	return cmd
}
