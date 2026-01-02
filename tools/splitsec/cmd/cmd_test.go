package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Kuniwak/ai-cli-tools/cli"
	"github.com/Kuniwak/ai-cli-tools/version"
	"github.com/google/go-cmp/cmp"
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

func TestMainCommandByArgs(t *testing.T) {
	testCases := map[string]struct {
		stdin    string
		args     []string
		expected map[string]string
	}{
		"empty": {
			stdin: "",
			args:  []string{"-tmpl", "section-%02d.md", "-out-dir", "%%OUTPUT_DIR%%"},
			expected: map[string]string{
				"section-00.md": "",
			},
		},
		"several sections": {
			stdin: "# section 1\n# section 2\n# section 3\n",
			args:  []string{"-tmpl", "section-%02d.md", "-out-dir", "%%OUTPUT_DIR%%"},
			expected: map[string]string{
				"section-00.md": "",
				"section-01.md": "# section 1\n",
				"section-02.md": "# section 2\n",
				"section-03.md": "# section 3\n",
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			outDir := t.TempDir()
			args := make([]string, len(tc.args))
			for i, arg := range tc.args {
				args[i] = strings.ReplaceAll(arg, "%%OUTPUT_DIR%%", outDir)
			}

			spy := cli.SpyProcInout(tc.stdin)
			exitStatus := MainCommandByArgs(args, spy.NewProcInout())

			if exitStatus != 0 {
				t.Logf("stderr: %s", spy.Stderr.String())
				t.Fatalf("expected exit status to be 0, got %d", exitStatus)
			}

			entries, err := os.ReadDir(outDir)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			actual := make(map[string]string)
			for _, entry := range entries {
				content, err := os.ReadFile(filepath.Join(outDir, entry.Name()))
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				actual[entry.Name()] = string(content)
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Error(cmp.Diff(tc.expected, actual))
			}
		})
	}
}
