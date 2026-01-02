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
	CommonOptions  tools.CommonOptions
	Reader         io.Reader
	CommandAndArgs []string
	Null           bool
	Parallel       int
}

func ParseOptions(args []string, inout *cli.ProcInout) (*Options, error) {
	flags := flag.NewFlagSet("stdinexec", flag.ContinueOnError)
	flags.SetOutput(inout.Stderr)
	flags.Usage = func() {
		fmt.Fprintf(inout.Stderr, `Usage: stdinexec [-0] [-p <parallel>] <command> [<args>...]

Execute a command for each line of the input, similar to "find -exec".

Options:
`)
		flags.PrintDefaults()
	}

	commonRawOptions := &tools.CommonRawOptions{}
	tools.DeclareCommonFlags(flags, commonRawOptions)

	null := flags.Bool("0", false, "use null byte as the record separator")
	parallelShort := flags.Int("p", 0, "number of parallel executions")
	parallelLong := flags.Int("parallel", 0, "number of parallel executions")

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

	var parallel int
	if *parallelLong != 0 {
		parallel = *parallelLong
	} else {
		parallel = *parallelShort
	}

	if parallel == 0 {
		parallel = 1
	} else if parallel < 0 {
		return nil, fmt.Errorf("ParseOptions: parallel must be at least 1")
	}

	if flags.NArg() == 0 {
		return nil, fmt.Errorf("ParseOptions: command is required")
	}

	commandAndArgs := flags.Args()

	return &Options{
		Reader:         inout.Stdin,
		CommandAndArgs: commandAndArgs,
		Null:           *null,
		Parallel:       parallel,
	}, nil
}
