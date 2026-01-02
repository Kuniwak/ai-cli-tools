package cmd

import (
	"fmt"

	"github.com/Kuniwak/ai-cli-tools/cli"
	"github.com/Kuniwak/ai-cli-tools/markdown"
	"github.com/Kuniwak/ai-cli-tools/version"
)

func MainCommandByArgs(args []string, inout *cli.ProcInout) int {
	options, err := ParseOptions(args, inout)
	if err != nil {
		fmt.Fprintln(inout.Stderr, err)
		return 1
	}
	if err := MainCommandByOptions(options, inout); err != nil {
		fmt.Fprintln(inout.Stderr, err)
		return 1
	}
	return 0
}

func MainCommandByOptions(options *Options, inout *cli.ProcInout) error {
	if options.CommonOptions.Help {
		return nil
	}

	if options.CommonOptions.Version {
		fmt.Fprintln(inout.Stdout, version.Version)
		return nil
	}

	if err := markdown.SplitBySections(options.Reader, markdown.NewWriterGenerator(options.OutputDirectory, options.BasenameTemplate)); err != nil {
		return fmt.Errorf("MainCommandByOptions: %w", err)
	}
	return nil
}
