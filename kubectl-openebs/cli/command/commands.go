package command

import (
	"github.com/openebs/openebsctl/kubectl-openebs/cli/command/cstor"
	"github.com/openebs/openebsctl/kubectl-openebs/cli/util"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewOpenebsCommand creates the `openebs` command and its nested children.
func NewOpenebsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubectl-openebs",
		Short: "OpenEBSctl is a a tool for interacting with OpenEBS storage components",
		Long:  `OpenEBSctl is a kubectl plugin to interact with OpenEBS container Attached Storage components. `,
	}
	kubernetesConfigFlags := genericclioptions.NewConfigFlags(true)
	kubernetesConfigFlags.AddFlags(cmd.PersistentFlags())

	cmd.AddCommand(
		util.NewCmdCompletion(cmd),
		cstor.NewCmdcStor(),
	)

	return cmd
}
