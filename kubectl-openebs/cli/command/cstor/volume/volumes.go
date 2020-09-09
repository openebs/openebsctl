/*
Copyright 2020 The OpenEBS Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package volume

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	volumeCommandHelpText = `# List Volumes:
	$ kubectl openebs cStor volume list

	# Info of a Volume:
	$ kubectl openebs cStor volume info --volname <vol>

  # Statistics of a Volume:
	$ kubectl openebs cStor volume stats --volname <vol>

  # Statistics of a Volume created in 'test' namespace:
	$ kubectl openebs cStor volume stats --volname <vol> --namespace test

  # Info of a Volume:
	$ kubectl openebs cStor volume describe --volname <vol>

  # Info of a Volume created in 'test' namespace:
	$ kubectl openebs cStor volume describe --volname <vol> --namespace test

  # Delete a Volume:
	$ kubectl openebs cStor volume delete --volname <vol>

  # Delete a Volume created in 'test' namespace:
	$ kubectl openebs cStor volume delete --volname <vol> --namespace test
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
		NewCmdVolumeInfo(),
		//TODO:
		//NewCmdVolumeStats(),
	)

	return cmd
}
