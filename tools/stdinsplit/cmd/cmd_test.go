package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Kuniwak/ai-cli-tools/cli"
	"github.com/Kuniwak/ai-cli-tools/version"
	"github.com/google/go-cmp/cmp"
)

func TestMainCommandByArgsHelp(t *testing.T) {
	spy := cli.SpyProcInout()

	exitStatus := MainCommandByArgs([]string{"-h"}, spy.NewProcInout())

	if exitStatus != 0 {
		t.Errorf("expected exit status to be 0, got %d", exitStatus)
	}

	if spy.Stderr.String() == "" {
		t.Error("expected stderr to be non-empty, got empty")
	}
}

func TestMainCommandByArgsVersion(t *testing.T) {
	spy := cli.SpyProcInout()
	exitStatus := MainCommandByArgs([]string{"-version"}, spy.NewProcInout())
	if exitStatus != 0 {
		t.Errorf("expected exit status to be 0, got %d", exitStatus)
	}
	expected := fmt.Sprintf("%s\n", version.Version)
	if spy.Stdout.String() != expected {
		t.Errorf("expected stdout to be %q, got %q", expected, spy.Stdout.String())
	}
}

func TestMainCommandByArgs(t *testing.T) {
	testCases := map[string]struct {
		stdin    string
		args     []string
		expected map[string]string
	}{
		"write to several files": {
			stdin:    "one\ntwo\n",
			args:     []string{"-l", "1", "-t", "test-%d.txt" /* -o t.TempDir() */},
			expected: map[string]string{"test-0.txt": "one\n", "test-1.txt": "two\n"},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			spy := cli.SpyProcInout(tc.stdin)
			tmpDir := t.TempDir()

			args := append(tc.args, "-o", tmpDir)
			exitStatus := MainCommandByArgs(args, spy.NewProcInout())
			if exitStatus != 0 {
				t.Errorf("expected exit status to be 0, got %d", exitStatus)
			}

			actual := make(map[string]string)
			entries, err := os.ReadDir(tmpDir)
			if err != nil {
				t.Errorf("expected error to be nil, got %v", err)
			}
			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}
				bs, err := os.ReadFile(filepath.Join(tmpDir, entry.Name()))
				if err != nil {
					t.Errorf("expected error to be nil, got %v", err)
				}
				actual[entry.Name()] = string(bs)
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Error(cmp.Diff(tc.expected, actual))
			}
		})
	}
}
