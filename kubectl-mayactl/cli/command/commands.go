package command

import (
	"github.com/spf13/cobra"

	"github.com/vaniisgh/mayactl/kubectl-mayactl/cli/command/cstor"
	"github.com/vaniisgh/mayactl/kubectl-mayactl/cli/util"
)

// NewMayaCommand creates the `mayactl` command and its nested children.
func NewMayaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubectl-mayactl",
		Short: "Maya means 'Magic' a tool for storage orchestration",
		Long:  `Maya means 'Magic' a tool for storage orchestration`,
	}

	cmd.AddCommand(
		util.NewCmdCompletion(cmd),
		cstor.NewCmdcStor(),
	)

	return cmd
}
