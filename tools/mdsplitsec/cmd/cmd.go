package cmd

import (
	"fmt"

	"github.com/Kuniwak/ai-cli-tools/cli"
	"github.com/Kuniwak/ai-cli-tools/lines"
	"github.com/Kuniwak/ai-cli-tools/markdown"
	"github.com/Kuniwak/ai-cli-tools/split"
	"github.com/Kuniwak/ai-cli-tools/testableio"
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

	outPathGenerator := split.NewOutPathgenerator(options.OutputDirectory, options.BasenameTemplate)
	openFileFunc := testableio.NewOpenFileFunc()
	filePaths, err := markdown.SplitBySections(options.Reader, outPathGenerator, openFileFunc)
	if err := lines.WriteLines(options.Null, filePaths, inout.Stdout); err != nil {
		return fmt.Errorf("MainCommandByOptions: failed to write file paths: %w", err)
	}
	if err != nil {
		return fmt.Errorf("MainCommandByOptions: %w", err)
	}
	return nil
}
