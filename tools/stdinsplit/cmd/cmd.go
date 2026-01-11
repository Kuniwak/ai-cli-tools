package cmd

import (
	"fmt"

	"github.com/Kuniwak/ai-cli-tools/cli"
	"github.com/Kuniwak/ai-cli-tools/lines"
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
	fmt.Fprintln(inout.Stderr, "Deprecated: stdinsplit is deprecated. Use GNU CoreUtils's split instead.")

	if options.CommonOptions.Help {
		return nil
	}

	if options.CommonOptions.Version {
		fmt.Fprintln(inout.Stdout, version.Version)
		return nil
	}

	outPathGenerator := split.NewOutPathgenerator(options.OutDir, options.Template)
	openFileFunc := testableio.NewOpenFileFunc()

	var writtenPaths []string
	var err error
	if options.LineCount != 0 {
		writtenPaths, err = split.SplitByLineCount(options.Reader, options.LineCount, outPathGenerator, openFileFunc)
		if err != nil {
			return fmt.Errorf("MainCommandByOptions: failed to split: %w", err)
		}
	} else if options.TotalCount != 0 {
		writtenPaths, err = split.SplitByTotalCount(options.Reader, options.TotalCount, outPathGenerator, openFileFunc)
		if err != nil {
			return fmt.Errorf("MainCommandByOptions: failed to split: %w", err)
		}
	} else {
		panic("either line count or total count must be specified")
	}

	if err := lines.WriteLines(options.Null, writtenPaths, inout.Stdout); err != nil {
		return fmt.Errorf("MainCommandByOptions: failed to write lines: %w", err)
	}

	return nil
}
