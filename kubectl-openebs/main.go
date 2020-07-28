package main

import (
	"github.com/openebs/openebsctl/kubectl-openebs/cli/command"
	"github.com/openebs/openebsctl/kubectl-openebs/cli/util"
)

func main() {
	err := command.NewMayaCommand().Execute()
	util.CheckError(err)
}
