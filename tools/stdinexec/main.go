package main

import (
	"github.com/Kuniwak/ai-cli-tools/cli"
	"github.com/Kuniwak/ai-cli-tools/tools/stdinexec/cmd"
)

func main() {
	cli.Run(cmd.MainCommandByArgs)
}
