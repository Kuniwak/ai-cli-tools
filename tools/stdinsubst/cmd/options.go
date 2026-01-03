package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/Kuniwak/ai-cli-tools/cli"
	"github.com/Kuniwak/ai-cli-tools/tools"
)

type Options struct {
	CommonOptions tools.CommonOptions
	Template      io.Reader
	Replacements  []Replacement
}

type Replacement struct {
	Before string
	After  string
}

func ParseOptions(args []string, inout *cli.ProcInout) (*Options, error) {
	flags := flag.NewFlagSet("stdinsubst", flag.ContinueOnError)
	flags.SetOutput(inout.Stderr)
	flags.Usage = func() {
		fmt.Fprintf(inout.Stderr, `Usage: stdinsubst <before-string> <after-file-path> [<before-string> <after-file-path> ...] < <template>

<before-string> and <after-file-path> are the strings to replace and the file path to read the replacement from.

Options:
`)
		flags.PrintDefaults()
	}

	commonRawOptions := &tools.CommonRawOptions{}
	tools.DeclareCommonFlags(flags, commonRawOptions)

	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return &Options{CommonOptions: tools.CommonOptions{Help: true}}, nil
		}
		return nil, fmt.Errorf("ParseOptions: %w", err)
	}

	commonOptions, err := tools.ValidateCommonOptions(commonRawOptions)
	if err != nil {
		return nil, fmt.Errorf("ParseOptions: %w", err)
	}

	if commonOptions.Version {
		return &Options{CommonOptions: commonOptions}, nil
	}

	if len(flags.Args())%2 != 0 {
		return nil, fmt.Errorf("number of replacements must be even")
	}

	replacements := make([]Replacement, len(flags.Args())/2)
	for i := 0; i < flags.NArg(); i += 2 {
		replacements[i/2] = Replacement{
			Before: flags.Arg(i),
			After:  flags.Arg(i + 1),
		}
	}

	return &Options{
		Replacements: replacements,
		Template:     inout.Stdin,
	}, nil
}
