package cmd

import (
	"fmt"
	"testing"

	"github.com/Kuniwak/ai-cli-tools/cli"
	"github.com/Kuniwak/ai-cli-tools/version"
)

func TestMainCommandByArgsHelp(t *testing.T) {
	spy := cli.SpyProcInout("")
	exitStatus := MainCommandByArgs([]string{"-h"}, spy.NewProcInout())
	if exitStatus != 0 {
		t.Errorf("expected exit status to be 0, got %d", exitStatus)
	}
	if spy.Stderr.String() == "" {
		t.Errorf("expected stderr to be non-empty, got %q", spy.Stderr.String())
	}
}

func TestMainCommandByArgsVersion(t *testing.T) {
	spy := cli.SpyProcInout("")
	exitStatus := MainCommandByArgs([]string{"-version"}, spy.NewProcInout())
	if exitStatus != 0 {
		t.Errorf("expected exit status to be 0, got %d", exitStatus)
	}
	expected := fmt.Sprintf("%s\n", version.Version)
	if spy.Stdout.String() != expected {
		t.Errorf("expected stdout to be %q, got %q", expected, spy.Stdout.String())
	}
}

func TestMainCommandByOptions(t *testing.T) {
	testCases := map[string]struct {
		stdin    string
		args     []string
		expected string
	}{
		"0 pairs": {
			stdin:    "%%GREETING%%, %%SUBJECT%%\n",
			args:     []string{},
			expected: "%%GREETING%%, %%SUBJECT%%\n",
		},
		"1 pair": {
			stdin:    "%%GREETING%%, %%NOUN%%!\n",
			args:     []string{"%%GREETING%%", "Hello"},
			expected: "Hello, %%NOUN%%!\n",
		},
		"several pairs": {
			stdin:    "%%GREETING%%, %%NOUN%%!\n%%GREETING%%, %%NOUN%%!\n",
			args:     []string{"%%GREETING%%", "Hello", "%%NOUN%%", "World"},
			expected: "Hello, World!\nHello, World!\n",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			spy := cli.SpyProcInout(tc.stdin)
			exitStatus := MainCommandByArgs(tc.args, spy.NewProcInout())
			if exitStatus != 0 {
				t.Errorf("expected exit status to be 0, got %d", exitStatus)
			}
			if spy.Stdout.String() != tc.expected {
				t.Errorf("expected stdout to be %q, got %q", tc.expected, spy.Stdout.String())
			}
		})
	}
}
