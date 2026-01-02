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
}

func ParseOptions(args []string, inout *cli.ProcInout) (*Options, error) {
	flags := flag.NewFlagSet("splitsec", flag.ContinueOnError)
	flags.SetOutput(inout.Stderr)
	flags.Usage = func() {
		fmt.Fprintf(inout.Stderr, `Usage: splitsec -i <input_file> -o <output_directory>
Split a markdown file by sections into files based on the number of seconds in each section.

Options:
`)
		flags.PrintDefaults()
	}

	outputDirectoryShort := flags.String("o", "", "output directory")
	outputDirectoryLong := flags.String("out-dir", "", "output directory")
	basenameTemplateShort := flags.String("t", "", "basename template")
	basenameTemplateLong := flags.String("tmpl", "", "basename template")

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
	}, nil
}
