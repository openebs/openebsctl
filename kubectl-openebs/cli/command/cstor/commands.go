package cstor

import (
	"flag"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/openebs/openebsctl/kubectl-openebs/cli/command/cstor/volume"
)

const (
	controllerStatusOk = "running"
	volumeStatusOK     = "Running"
	timeout            = 5 * time.Second

	//CASType defines the Sotrage engine used
	CASType = "cStor"
)

var (
	cStorCommandHelpText = `
The following commands helps in retreiving information of the cStor realted to
volumes, pools,  and so on.

Usage: kubectl openebs cStor <subcommand> [options] [args]

Examples:

 # Status
	$ kubectl openebs cStor status

 #Volume
	$ kubectl openebs cStor volume --help

`
	namespace string
)

// NewCmdcStor provides options for managing OpenEBS Volume
func NewCmdcStor() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "cStor",
		Aliases: []string{"cstor"},
		Short:   "Provides operations related to a cStor storage engine",
		Long:    cStorCommandHelpText,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(cStorCommandHelpText)
		},
	}

	cmd.AddCommand(

		volume.NewCmdVolume(cmd),
		//TODO: uncomment all one by one
		//NewCmdVolumeDelete(),
		//NewCmdVolumeStats(),
		//NewCmdVolumeInfo(),
	)

	cmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "openebs",
		"namespace name, required if volume is not in the `default` namespace")

	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	flag.CommandLine.Parse([]string{})
	viper.BindPFlag("namespace", cmd.PersistentFlags().Lookup("namespace"))

	return cmd
}
