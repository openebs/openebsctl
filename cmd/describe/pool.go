package describe

import (
	"github.com/openebs/openebsctl/pkg/cstor/describe"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/spf13/cobra"
)

var (
	poolInfoCommandHelpText = `
This command fetches information and status of the various aspects 
of the cStor Pool Instance and its underlying related resources in the provided namespace.
If no namespace is provided it uses default namespace for execution.
$ kubectl openebs describe pool [cspi-name] -n [namespace]
`
)

// NewCmdDescribePool displays OpenEBS cStor pool instance information.
func NewCmdDescribePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pool",
		Aliases: []string{"pools", "p"},
		Short:   "Displays cStorPoolInstance information",
		Long:    poolInfoCommandHelpText,
		Example: `kubectl openebs describe pool cspi-one -n openebs`,
		Run: func(cmd *cobra.Command, args []string) {
			var namespace string // This namespace belongs to the CSPI entered
			if namespace, _ = cmd.Flags().GetString("namespace"); namespace == "" {
				// NOTE: The error comes as nil even when the ns flag is not specified
				namespace = "openebs"
			}
			util.CheckErr(describe.RunPoolInfo(cmd, args, namespace), util.Fatal)
		},
	}
	return cmd
}
