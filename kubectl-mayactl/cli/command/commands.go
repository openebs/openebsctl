package command

import (
	"flag"

	"github.com/spf13/cobra"

	"github.com/vaniisgh/mayactl/kubectl-mayactl/cli/command/volume"
	"github.com/vaniisgh/mayactl/kubectl-mayactl/cli/util"
)

var test string
var namespace string

// NewMayaCommand creates the `mayactl` command and its nested children.
func NewMayaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubectl-mayactl",
		Short: "Maya means 'Magic' a tool for storage orchestration",
		Long:  `Maya means 'Magic' a tool for storage orchestration`,
	}

	cmd.AddCommand(
		util.NewCmdCompletion(cmd),
		volume.NewCmdVolume(cmd),
	)

	cmd.PersistentFlags().StringVarP(&test, "test", "t", "this is a test", "")

	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	flag.CommandLine.Parse([]string{})

	//fmt.Println("here" + test)
	return cmd
}
