package describe

import (
	"github.com/openebs/openebsctl/pkg/cstor/describe"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/spf13/cobra"
)

var (
	volumeInfoCommandHelpText = `
This command fetches information and status of the various
aspects of a cStor Volume such as ISCSI, Controller, and Replica.

#
$ kubectl openebs describe [pool|volume] [name]

`
)

// NewCmdDescribeVolume displays OpenEBS Volume information.
func NewCmdDescribeVolume() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "volume",
		Aliases: []string{"volumes", "vol", "v"},
		Short:   "Displays Openebs information",
		Long:    volumeInfoCommandHelpText,
		Example: `kubectl openebs describe volume [vol]`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Get this from flags, pflag, etc
			var ns string
			if ns, _ = cmd.Flags().GetString("namespace"); ns == "" {
				// NOTE: The error comes as nil even when the ns flag is not specified
				ns = "openebs"
			}
			util.CheckErr(describe.RunVolumeInfo(cmd, args, ns), util.Fatal)
		},
	}
	return cmd
}
