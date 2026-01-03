package main

import (
	"github.com/Kuniwak/ai-cli-tools/cli"
	"github.com/Kuniwak/ai-cli-tools/tools/mdsplitsec/cmd"
)

func main() {
	cli.Run(cmd.MainCommandByArgs)
}
