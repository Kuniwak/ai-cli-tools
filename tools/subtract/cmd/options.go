package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Kuniwak/ai-cli-tools/cli"
	"github.com/Kuniwak/ai-cli-tools/tools"
)

type Options struct {
	CommonOptions tools.CommonOptions
	Null          bool
	Minuend       io.Reader
	Subtrahends   []io.ReadCloser
}

func ParseOptions(args []string, inout *cli.ProcInout) (*Options, error) {
	flags := flag.NewFlagSet("subtract", flag.ContinueOnError)
	flags.SetOutput(inout.Stderr)
	flags.Usage = func() {
		fmt.Fprintf(inout.Stderr, `Usage: subtract <number1> <number2>

Subtract the second number from the first number.

Options:
`)
		flags.PrintDefaults()

		fmt.Fprintf(inout.Stderr, `
Examples:
  $ cat ./minuend.txt
  line 1
  line 2
  line 3

  $ cat ./subtrahend1.txt
  line 2

  $ cat ./subtrahend2.txt
  line 2
  line 3

  $ subtract ./subtrahend1.txt ./subtrahend2.txt < ./minuend.txt 
  line 1

  $ # It is useful to drop processed files from the input.
  $ find ./input -name '*.md' -print0 | subtract -0 <(find ./output -name '*.md' -print0 | sed -e 's|^\./input/|./output/|')
  ./input/file1.md
  ...

  $ # Process unprocessed ./input/*.md files in parallel using 3 processes by Claude Code.
  $ subtract -0 <(find ./input -name '*.md' -print0) <(find ./output -name '*.md' -print0 | sed -e 's|^\./input/|./output/|') | stdinexec -0 bash -c 'claude -p < "{}"'
`)
	}

	commonRawOptions := &tools.CommonRawOptions{}
	tools.DeclareCommonFlags(flags, commonRawOptions)

	null := flags.Bool("0", false, "use null byte as the record separator")

	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return &Options{CommonOptions: tools.CommonOptions{Help: true}}, nil
		}
		return nil, fmt.Errorf("ParseOptions: failed to parse options: %w", err)
	}

	commonOptions, err := tools.ValidateCommonOptions(commonRawOptions)
	if err != nil {
		return nil, fmt.Errorf("ParseOptions: failed to validate common options: %w", err)
	}

	if commonOptions.Version {
		return &Options{CommonOptions: commonOptions}, nil
	}

	subtrahends := make([]io.ReadCloser, flags.NArg())
	for i := 0; i < flags.NArg(); i++ {
		subtrahend, err := os.OpenFile(flags.Arg(i), os.O_RDONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("ParseOptions: failed to open subtrahend file: %w at %q", err, flags.Arg(i))
		}
		subtrahends[i] = subtrahend
	}

	return &Options{
		CommonOptions: commonOptions,
		Null:          *null,
		Minuend:       inout.Stdin,
		Subtrahends:   subtrahends,
	}, nil
}
