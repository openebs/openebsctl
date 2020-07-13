package volume

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	volumeCommandHelpText = `# List Volumes:
	$ kubectl mayactl cStor volume list

  # Statistics of a Volume:
	$ kubectl mayactl cStor volume stats --volname <vol>

  # Statistics of a Volume created in 'test' namespace:
	$ kubectl mayactl cStor volume stats --volname <vol> --namespace test

  # Info of a Volume:
	$ kubectl mayactl cStor volume describe --volname <vol>

  # Info of a Volume created in 'test' namespace:
	$ kubectl mayactl cStor volume describe --volname <vol> --namespace test

  # Delete a Volume:
	$ kubectl mayactl cStor volume delete --volname <vol>

  # Delete a Volume created in 'test' namespace:
	$ kubectl mayactl cStor volume delete --volname <vol> --namespace test
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
