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
	CommonOptions    tools.CommonOptions
	Reader           io.Reader
	OutputDirectory  string
	BasenameTemplate string
	Null             bool
}

func ParseOptions(args []string, inout *cli.ProcInout) (*Options, error) {
	flags := flag.NewFlagSet("mdsplitsec", flag.ContinueOnError)
	flags.SetOutput(inout.Stderr)
	flags.Usage = func() {
		fmt.Fprintf(inout.Stderr, `Usage: mdsplitsec [-0] -o <output_directory> < <markdown>

Split a markdown file by sections into files based on the number of seconds in each section.

Options:
`)
		flags.PrintDefaults()
		fmt.Fprintf(inout.Stderr, `
Examples:
  $ cat ./input.md
  # Section 1
  ...
  # Section 2
  ...
  # Section 3
  ...

  $ mdsplitsec -o ./output < ./input.md
  ./output/section-01.md
  ./output/section-02.md
  ./output/section-03.md

  $ cat ./output/section-01.md
  # Section 1
  ...
`)
	}

	outputDirectoryShort := flags.String("o", "", "output directory")
	outputDirectoryLong := flags.String("out-dir", "", "output directory")
	basenameTemplateShort := flags.String("t", "", "basename template")
	basenameTemplateLong := flags.String("tmpl", "", "basename template")
	null := flags.Bool("0", false, "use null byte as the record separator")

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

	var outputDirectory string
	if *outputDirectoryLong != "" {
		outputDirectory = *outputDirectoryLong
	} else {
		outputDirectory = *outputDirectoryShort
	}

	var basenameTemplate string
	if *basenameTemplateLong != "" {
		basenameTemplate = *basenameTemplateLong
	} else {
		basenameTemplate = *basenameTemplateShort
	}

	if outputDirectory == "" {
		return nil, fmt.Errorf("ParseOptions: output directory is required")
	}

	stat, err := os.Stat(outputDirectory)
	if err == nil {
		if !stat.IsDir() {
			return nil, fmt.Errorf("ParseOptions: output directory is not a directory: %s", outputDirectory)
		}
	} else {
		if os.IsNotExist(err) {
			os.MkdirAll(outputDirectory, 0755)
		} else {
			return nil, fmt.Errorf("ParseOptions: failed to create output directory: %s", outputDirectory)
		}
	}

	if basenameTemplate == "" {
		basenameTemplate = "section-%02d.md"
	}

	return &Options{
		Reader:           inout.Stdin,
		OutputDirectory:  outputDirectory,
		BasenameTemplate: basenameTemplate,
		Null:             *null,
	}, nil
}
