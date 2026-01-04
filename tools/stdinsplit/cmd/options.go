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
	Reader        io.Reader
	Null          bool
	OutDir        string
	Template      string
	LineCount     int
	TotalCount    int
}

func ParseOptions(args []string, inout *cli.ProcInout) (*Options, error) {
	flags := flag.NewFlagSet("stdinsplit", flag.ContinueOnError)
	flags.SetOutput(inout.Stderr)
	flags.Usage = func() {
		fmt.Fprintf(inout.Stderr, `Usage: stdinsplit [-0] (-l <line-count> | -n <total-count>) -o <out-dir> [-t <template>] < <input>

Split the input by the separator and write each part to a file in the output directory.
If <line-count> is specified, split the input into <line-count> lines.
If <total-count> is specified, split the input into <total-count> parts.

Options:
`)
		flags.PrintDefaults()

		fmt.Fprintln(inout.Stderr, `
Examples:
  $ # Split the input into 10 parts.
  $ echo "Hello\nWorld\n" | stdinsplit -o ./output -n 10 -t "part-%02d.txt"
  ./output/part-00.txt
  ./output/part-01.txt
  ...
  ./output/part-09.txt

  $ # Split the input into 10 lines.
  $ echo "Hello\nWorld\n" | stdinsplit -o ./output -l 1 -t "part-%02d.txt"
  ./output/part-00.txt
  ./output/part-01.txt

  $ # Use with stdinexec to process each part in parallel.
  $ echo "Hello\nWorld\n" | stdinsplit -0 -o ./output -l 1 | stdinexec -0 -p 2 bash -c 'claude -p < "{}"'`)
	}

	commonRawOptions := &tools.CommonRawOptions{}
	tools.DeclareCommonFlags(flags, commonRawOptions)

	null := flags.Bool("0", false, "use null byte as the record separator")
	templateLong := flags.String("t", "", "basename template (default: \"%03d.txt\")")
	templateShort := flags.String("template", "", "basename template (default: \"%03d.txt\")")
	outDirLong := flags.String("out-dir", "", "output directory")
	outDirShort := flags.String("o", "", "output directory h")
	lineCountShort := flags.Int("l", 0, "number of lines per part")
	lineCountLong := flags.Int("line-count", 0, "number of lines per part")
	totalCountShort := flags.Int("n", 0, "number of parts")
	totalCountLong := flags.Int("total-count", 0, "number of parts")

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

	var outDir string
	if *outDirLong != "" {
		outDir = *outDirLong
	} else {
		outDir = *outDirShort
	}

	if outDir == "" {
		return nil, fmt.Errorf("ParseOptions: output directory is required")
	}

	if stat, err := os.Stat(outDir); err == nil {
		if !stat.IsDir() {
			return nil, fmt.Errorf("ParseOptions: output directory is not a directory")
		}
	} else if !stat.IsDir() {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(outDir, 0755); err != nil {
				return nil, fmt.Errorf("ParseOptions: failed to create output directory: %w", err)
			}
		} else {
			return nil, fmt.Errorf("ParseOptions: failed to stat output directory: %w", err)
		}
	}

	var lineCount int
	if *lineCountLong != 0 {
		lineCount = *lineCountLong
	} else {
		lineCount = *lineCountShort
	}

	var totalCount int
	if *totalCountLong != 0 {
		totalCount = *totalCountLong
	} else {
		totalCount = *totalCountShort
	}

	if lineCount == 0 && totalCount == 0 {
		return nil, fmt.Errorf("ParseOptions: either line-count or total-count must be specified")
	}

	if lineCount != 0 && totalCount != 0 {
		return nil, fmt.Errorf("ParseOptions: either line-count or total-count must be specified")
	}

	var template string
	if *templateLong != "" {
		template = *templateLong
	} else {
		template = *templateShort
	}

	if template == "" {
		template = "%03d.txt"
	}

	return &Options{
		CommonOptions: commonOptions,
		Reader:        inout.Stdin,
		Null:          *null,
		OutDir:        outDir,
		Template:      template,
		LineCount:     lineCount,
		TotalCount:    totalCount,
	}, nil
}
