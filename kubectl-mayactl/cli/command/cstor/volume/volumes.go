package volume

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	volumeCommandHelpText = `# List Volumes:
	$ kubectl mayactl jiba volume list

  # Statistics of a Volume:
	$ kubectl mayactl zfs volume stats --volname <vol>

  #TODO: fix commands here
  # Statistics of a Volume created in 'test' namespace:
	$ mayactl volume stats --volname <vol> --namespace test

  # Info of a Volume:
	$ mayactl volume describe --volname <vol>

  # Info of a Volume created in 'test' namespace:
	$ mayactl volume describe --volname <vol> --namespace test

  # Delete a Volume:
	$ mayactl volume delete --volname <vol>

  # Delete a Volume created in 'test' namespace:
	$ mayactl volume delete --volname <vol> --namespace test
 `
)

// NewCmdVolume provides options for managing OpenEBS Volume
func NewCmdVolume(rootCmd *cobra.Command) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "volume",
		Short: "Provides operations related to a Volume",
		Long:  volumeCommandHelpText,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(volumeCommandHelpText)
		},
	}

	cmd.AddCommand(
		NewCmdVolumesList(),
		//NewCmdVolumesInfo(),
	)

	return cmd
}
