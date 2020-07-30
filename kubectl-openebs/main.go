package main

import (
	"github.com/openebs/openebsctl/kubectl-openebs/cli/command"
	"github.com/openebs/openebsctl/kubectl-openebs/cli/util"
)

func main() {
	err := command.NewOpenebsCommand().Execute()
	util.CheckError(err)
}
