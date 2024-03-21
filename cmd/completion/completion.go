/*
Copyright 2020-2022 The OpenEBS Authors

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

package completion

import (
	"io"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

const (
	completionUsage = `
To load completion to current bash shell,
. <(kubectl openebs completion bash)

To configure your bash shell to load completions for each session add to your bashrc
# ~/.bashrc or ~/.profile
. <(kubectl openebs completion bash)

To load completion to current zsh shell,
. <(kubectl openebs completion zsh)

To configure your zsh shell to load completions for each session add to your zshrc
# ~/.zshrc
. <(kubectl openebs completion zsh)

Do similar steps for fish & powershell
`
)

// NewCmdCompletion creates the completion command
func NewCmdCompletion(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:       "completion",
		ValidArgs: []string{"bash", "zsh"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Short:     "Outputs shell completion code for the specified shell (bash or zsh)",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
		Run: func(cmd *cobra.Command, args []string) {
			RunCompletion(os.Stdout, rootCmd, args)
		},
	}
	cmd.SetUsageTemplate(completionUsage)
	return cmd
}

// RunCompletion is used to run the completion of the cobra commad
func RunCompletion(out io.Writer, cmd *cobra.Command, args []string) {
	var err error
	switch args[0] {
	case "bash":
		err = cmd.GenBashCompletion(out)
	case "zsh":
		err = cmd.GenZshCompletion(out)
	case "fish":
		err = cmd.GenFishCompletion(out, true)
	case "powershell":
		err = cmd.GenPowerShellCompletion(out)
	}
	if err != nil {
		klog.Error(err)
	}
}
