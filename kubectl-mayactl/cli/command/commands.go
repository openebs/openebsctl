package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vaniisgh/mayactl/kubectl-mayactl/cli/utils"
)

// NewMayaCommand creates the `maya` command and its nested children.
func NewMayaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubectl-mayactl",
		Short: "Maya means 'Magic' a tool for storage orchestration",
		Long:  `Maya means 'Magic' a tool for storage orchestration`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			fmt.Println("command creating")
		},
	}

	cmd.AddCommand(
		utils.NewCmdCompletion(cmd),
	)

	//TODO: declare in k8s.go
	//cmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "openebs",
	//	"namespace name, required if volume is not in the default namespace")

	//flag.CommandLine.Parse([]string{})

	return cmd
}
