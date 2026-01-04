package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"slices"
	"strconv"
	"strings"

	"github.com/Kuniwak/ai-cli-tools/cli"
	"github.com/Kuniwak/ai-cli-tools/lines"
	"github.com/Kuniwak/ai-cli-tools/version"
	"golang.org/x/sync/errgroup"
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

	ch := make(chan string)
	var eg errgroup.Group
	eg.Go(func() error {
		defer close(ch)

		scanner := bufio.NewScanner(options.Reader)
		scanFunc := lines.NewScanFunc(options.Null)
		scanner.Split(scanFunc)
		for scanner.Scan() {
			ch <- scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("MainCommandByOptions: failed to scan lines: %w", err)
		}

		return nil
	})

	for i := 0; i < options.Parallel; i++ {
		eg.Go(executeCommand(i, options.Parallel, ch, slices.Clone(options.CommandAndArgs), inout))
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("MainCommandByOptions: failed to wait for commands to complete: %w", err)
	}
	return nil
}

func executeCommand(i int, parallel int, lines <-chan string, commandAndArgsOrig []string, inout *cli.ProcInout) func() error {
	return func() error {
		for line := range lines {
			commandAndArgs := slices.Clone(commandAndArgsOrig)
			for j, arg := range commandAndArgs {
				commandAndArgs[j] = strings.ReplaceAll(arg, "{}", line)
			}
			cmd := exec.Command(commandAndArgs[0], commandAndArgs[1:]...)

			stdoutReader, _ := cmd.StdoutPipe()
			stderrReader, _ := cmd.StderrPipe()

			if err := cmd.Start(); err != nil {
				return fmt.Errorf("MainCommandByOptions: failed to execute command: %w (%d %#v)", err, i, commandAndArgs)
			}

			var eg errgroup.Group
			eg.Go(func() error {
				scanner := bufio.NewScanner(stdoutReader)
				scanner.Split(bufio.ScanLines)
				for scanner.Scan() {
					text := scanner.Text()
					if parallel == 1 {
						fmt.Fprintln(inout.Stdout, text)
					} else {
						io.WriteString(inout.Stdout, strconv.Itoa(i))
						io.WriteString(inout.Stdout, "\t")
						io.WriteString(inout.Stdout, text)
						io.WriteString(inout.Stdout, "\n")
					}
				}
				if err := scanner.Err(); err != nil {
					return fmt.Errorf("MainCommandByOptions: failed to scan stdout: %w", err)
				}
				return nil
			})

			eg.Go(func() error {
				scanner := bufio.NewScanner(stderrReader)
				scanner.Split(bufio.ScanLines)
				for scanner.Scan() {
					text := scanner.Text()
					if parallel == 1 {
						fmt.Fprintln(inout.Stdout, text)
					} else {
						io.WriteString(inout.Stdout, strconv.Itoa(i))
						io.WriteString(inout.Stdout, "\t")
						io.WriteString(inout.Stdout, text)
						io.WriteString(inout.Stdout, "\n")
					}
				}
				if err := scanner.Err(); err != nil {
					return fmt.Errorf("MainCommandByOptions: failed to scan stderr: %w", err)
				}
				return nil
			})

			if err := eg.Wait(); err != nil {
				return fmt.Errorf("MainCommandByOptions: failed to wait for stdout and stderr to complete: %w", err)
			}

			if err := cmd.Wait(); err != nil {
				return fmt.Errorf("MainCommandByOptions: failed to wait for command to complete: %w (%d %#v)", err, i, commandAndArgs)
			}
		}
		return nil
	}
}
