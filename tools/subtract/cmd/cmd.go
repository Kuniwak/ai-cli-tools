package cmd

import (
	"bufio"
	"fmt"
	"io"

	"github.com/Kuniwak/ai-cli-tools/cli"
	"github.com/Kuniwak/ai-cli-tools/scanners"
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

	var splitFunc bufio.SplitFunc
	if options.Null {
		splitFunc = scanners.ScanLinesWithNull
	} else {
		splitFunc = bufio.ScanLines
	}

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
			io.WriteString(inout.Stdout, minuend)
			if options.Null {
				io.WriteString(inout.Stdout, "\u0000")
			} else {
				io.WriteString(inout.Stdout, "\n")
			}
		}
	}

	if err := minuendScanner.Err(); err != nil {
		return fmt.Errorf("MainCommandByOptions: failed to scan minuend: %w", err)
	}

	return nil
}
