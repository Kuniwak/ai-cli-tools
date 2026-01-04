package cmd

import (
	"bufio"
	"fmt"

	"github.com/Kuniwak/ai-cli-tools/cli"
	"github.com/Kuniwak/ai-cli-tools/lines"
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

	splitFunc := lines.NewScanFunc(options.Null)

	m := make(map[string]struct{})
	for _, subtrahend := range options.Subtrahends {
		subtrahendScanner := bufio.NewScanner(subtrahend)
		subtrahendScanner.Split(splitFunc)
		for subtrahendScanner.Scan() {
			m[subtrahendScanner.Text()] = struct{}{}
		}
		if err := subtrahendScanner.Err(); err != nil {
			return fmt.Errorf("MainCommandByOptions: failed to scan subtrahend: %w", err)
		}
		_ = subtrahend.Close()
	}

	minuendScanner := bufio.NewScanner(options.Minuend)
	minuendScanner.Split(splitFunc)
	for minuendScanner.Scan() {
		minuend := minuendScanner.Text()
		if _, ok := m[minuend]; !ok {
			if err := lines.WriteLine(options.Null, minuend, inout.Stdout); err != nil {
				return fmt.Errorf("MainCommandByOptions: failed to write minuend: %w", err)
			}
		}
	}

	if err := minuendScanner.Err(); err != nil {
		return fmt.Errorf("MainCommandByOptions: failed to scan minuend: %w", err)
	}

	return nil
}
