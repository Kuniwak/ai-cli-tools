package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Kuniwak/ai-cli-tools/cli"
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
	if spy.Stdout.String() == "" {
		t.Errorf("expected stdout to be non-empty, got %q", spy.Stdout.String())
	}
}

func TestMainCommandByArgs(t *testing.T) {

	testCases := map[string]struct {
		stdin    string
		null     bool
		files    map[string]string
		expected string
	}{
		"empty": {
			stdin:    "",
			files:    map[string]string{},
			expected: "",
		},
		"one line": {
			stdin:    "line 1\n",
			files:    map[string]string{},
			expected: "line 1\n",
		},
		"one line in the same subtrahend as the minuend": {
			stdin:    "line 1\n",
			files:    map[string]string{"1": "line 1\n"},
			expected: "",
		},
		"one line in different subtrahend": {
			stdin:    "line 1\n",
			files:    map[string]string{"2": "line 2\n"},
			expected: "line 1\n",
		},
		"several files": {
			stdin:    "line 1\nline 2\nline 3\n",
			files:    map[string]string{"2": "line 2\n", "3": "line 3\n"},
			expected: "line 1\n",
		},
		"null": {
			stdin:    "line 1\u0000line 2\u0000line 3\u0000",
			null:     true,
			files:    map[string]string{"1": "line 1\u0000"},
			expected: "line 2\u0000line 3\u0000",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tmpDir := t.TempDir()
			args := make([]string, 0, len(tc.files)+1)
			if tc.null {
				args = append(args, "-0")
			}
			for basename, content := range tc.files {
				filePath := filepath.Join(tmpDir, basename)
				if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				args = append(args, filePath)
			}

			spy := cli.SpyProcInout(tc.stdin)
			exitStatus := MainCommandByArgs(args, spy.NewProcInout())
			if exitStatus != 0 {
				t.Fatalf("expected exit status to be 0, got %d\n%s", exitStatus, spy.Stderr.String())
			}
			if spy.Stdout.String() != tc.expected {
				t.Errorf("expected stdout to be %q, got %q", tc.expected, spy.Stdout.String())
			}
		})
	}
}
