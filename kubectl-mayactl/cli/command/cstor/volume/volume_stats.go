package volume

import (

	//"github.com/openebs/api/pkg/apis/openebs.io/v1alpha1"
	//v1 "github.com/openebs/api/types/v1"

	"github.com/spf13/cobra"
)

var (
	volumeStatsCommandHelpText = `
This command queries the statisics of a volume.

Usage: mayactl volume stats --volname <vol> [-size <size>]
`
)

// NewCmdVolumeStats displays the runtime statistics of volume
func NewCmdVolumeStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stats",
		Short:   "Displays the runtime statisics of Volume",
		Long:    volumeStatsCommandHelpText,
		Example: ` mayactl volume stats --volname=vol`,
		Run: func(cmd *cobra.Command, args []string) {
			runVolumeStats(cmd)
			//util.CheckErr(options.runVolumeStats(cmd), util.Fatal)
		},
	}

	cmd.Flags().StringVarP(&volName, "volname", "", volName,
		"unique volume name.")
	return cmd
}

func runVolumeStats(cmd *cobra.Command) error {

	return nil
}
