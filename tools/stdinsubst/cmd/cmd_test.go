package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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
		args     []StringGenerator
		expected string
	}{
		"0 pairs": {
			stdin:    "%%GREETING%%, %%SUBJECT%%\n",
			args:     []StringGenerator{},
			expected: "%%GREETING%%, %%SUBJECT%%\n",
		},
		"1 pair": {
			stdin:    "%%GREETING%%, %%NOUN%%!\n",
			args:     []StringGenerator{constantString("%%GREETING%%"), filePath("hello.txt", "Hello")},
			expected: "Hello, %%NOUN%%!\n",
		},
		"several pairs": {
			stdin:    "%%GREETING%%, %%NOUN%%!\n",
			args:     []StringGenerator{constantString("%%GREETING%%"), filePath("hello.txt", "Hello"), constantString("%%NOUN%%"), filePath("world.txt", "World")},
			expected: "Hello, World!\n",
		},
		"several occurrences": {
			stdin:    "%%GREETING%%\n%%GREETING%%\n",
			args:     []StringGenerator{constantString("%%GREETING%%"), filePath("hello.txt", "Hello")},
			expected: "Hello\nHello\n",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			spy := cli.SpyProcInout(tc.stdin)
			exitStatus := MainCommandByArgs(renderString(tc.args, t.TempDir()), spy.NewProcInout())
			if exitStatus != 0 {
				t.Errorf("expected exit status to be 0, got %d\n%s", exitStatus, spy.Stderr.String())
			}
			if spy.Stdout.String() != tc.expected {
				t.Errorf("expected stdout to be %q, got %q", tc.expected, spy.Stdout.String())
			}
		})
	}
}

type StringGenerator func(outDir string) string

func constantString(s string) StringGenerator {
	return func(_ string) string {
		return s
	}
}

func filePath(basename string, content string) StringGenerator {
	return func(outDir string) string {
		filePath := filepath.Join(outDir, basename)
		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		if _, err := io.WriteString(f, content); err != nil {
			panic(err)
		}
		return filePath
	}
}

func renderString(gs []StringGenerator, outDir string) []string {
	rendered := make([]string, len(gs))
	for i, g := range gs {
		rendered[i] = g(outDir)
	}
	return rendered
}
