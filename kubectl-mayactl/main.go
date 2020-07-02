package main

import (
	"github.com/vaniisgh/mayactl/kubectl-mayactl/cli/command"
	"github.com/vaniisgh/mayactl/kubectl-mayactl/cli/util"
)

func main() {
	err := command.NewMayaCommand().Execute()
	util.CheckError(err)
}
