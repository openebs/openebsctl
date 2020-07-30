package command

import (
	"github.com/spf13/cobra"

	"github.com/openebs/openebsctl/kubectl-openebs/cli/command/cstor"
	"github.com/openebs/openebsctl/kubectl-openebs/cli/util"
)

// NewOpenebsCommand creates the `openebs` command and its nested children.
func NewOpenebsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubectl-openebs",
		Short: "OpenEBSctl is a a tool for interacting with OpenEBS storage components",
		Long:  `OpenEBSctl is a kubectl plugin to interact with OpenEBS container Attached Storage components. `,
	}

	cmd.AddCommand(
		util.NewCmdCompletion(cmd),
		cstor.NewCmdcStor(),
	)

	return cmd
}
