package cmd

import (
	"fmt"
	"io"

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

	outPathGenerator := markdown.NewOutPathgenerator(options.OutputDirectory, options.BasenameTemplate)
	openFileFunc := markdown.NewOpenFileFunc()
	filePaths, err := markdown.SplitBySections(options.Reader, outPathGenerator, openFileFunc)
	if err := WriteFilePaths(options.Null, filePaths, inout); err != nil {
		return fmt.Errorf("MainCommandByOptions: failed to write file paths: %w", err)
	}
	if err != nil {
		return fmt.Errorf("MainCommandByOptions: %w", err)
	}
	return nil
}

func WriteFilePaths(null bool, filePaths []string, inout *cli.ProcInout) error {
	for _, filePath := range filePaths {
		if null {
			io.WriteString(inout.Stdout, filePath)
			io.WriteString(inout.Stdout, "\000")
		} else {
			fmt.Fprintln(inout.Stdout, filePath)
		}
	}
	return nil
}
