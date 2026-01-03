package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Kuniwak/ai-cli-tools/cli"
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

	args := make([]string, len(options.Replacements)*2)
	for i, replacement := range options.Replacements {
		args[i*2] = replacement.Before
		bs, err := os.ReadFile(replacement.After)
		if err != nil {
			return fmt.Errorf("MainCommandByOptions: failed to read replacement file: %w", err)
		}
		args[i*2+1] = string(bs)
	}
	replacer := strings.NewReplacer(args...)
	tmpl, err := io.ReadAll(options.Template)
	if err != nil {
		return fmt.Errorf("MainCommandByOptions: failed to read template: %w", err)
	}
	io.WriteString(inout.Stdout, replacer.Replace(string(tmpl)))
	return nil
}
