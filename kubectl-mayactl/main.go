package main

import (
	"github.com/vaniisgh/mayactl/kubectl-mayactl/cli/command"
	"github.com/vaniisgh/mayactl/kubectl-mayactl/cli/utils"
)

func main() {
	err := command.NewMayaCommand().Execute()

	utils.CheckError(err)
}
