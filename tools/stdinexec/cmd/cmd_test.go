package cmd

import (
	"fmt"
	"slices"
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

func TestMainCommandByArgsSingle(t *testing.T) {
	testCases := map[string]struct {
		stdin          string
		args           []string
		expectedStdout string
	}{
		"empty": {
			stdin:          "",
			args:           []string{"echo", "hello", "{}"},
			expectedStdout: "",
		},
		"one line": {
			stdin:          "one\n",
			args:           []string{"echo", "hello", "{}"},
			expectedStdout: "hello one\n",
		},
		"several lines": {
			stdin:          "one\ntwo\nthree\n",
			args:           []string{"echo", "hello", "{}"},
			expectedStdout: "hello one\nhello two\nhello three\n",
		},
		"null": {
			stdin:          "one\u0000two\u0000three\u0000",
			args:           []string{"-0", "echo", "hello", "{}"},
			expectedStdout: "hello one\nhello two\nhello three\n",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			spy := cli.SpyProcInout(tc.stdin)
			exitStatus := MainCommandByArgs(tc.args, spy.NewProcInout())
			if exitStatus != 0 {
				t.Errorf("expected exit status to be 0, got %d\n%s", exitStatus, spy.Stderr.String())
			}
			if spy.Stdout.String() != tc.expectedStdout {
				t.Error(cmp.Diff(tc.expectedStdout, spy.Stdout.String()))
			}
		})
	}
}

func TestMainCommandByArgsParallel(t *testing.T) {
	expectedStdoutOneOf := make([]string, 0)
	var msg [4]string
	for _, processedBy0_1 := range []bool{false, true} {
		for _, processedBy0_2 := range []bool{false, true} {
			for _, oneTwo := range []bool{false, true} {
				if processedBy0_1 {
					msg[0] = "0"
				} else {
					msg[0] = "1"
				}

				if processedBy0_2 {
					msg[2] = "0"
				} else {
					msg[2] = "1"
				}

				if oneTwo {
					msg[1] = "hello one"
					msg[3] = "hello two"
				} else {
					msg[1] = "hello two"
					msg[3] = "hello one"
				}

				expectedStdoutOneOf = append(expectedStdoutOneOf, fmt.Sprintf("%s\t%s\n%s\t%s\n", msg[0], msg[1], msg[2], msg[3]))
			}
		}
	}

	for i := 1; i <= 10; i++ {
		t.Run(fmt.Sprintf("parallel %d", i), func(t *testing.T) {
			spy := cli.SpyProcInout("one\ntwo\n")
			exitStatus := MainCommandByArgs([]string{"-p", "2", "echo", "hello", "{}"}, spy.NewProcInout())
			if exitStatus != 0 {
				t.Errorf("expected exit status to be 0, got %d", exitStatus)
			}

			if !slices.Contains(expectedStdoutOneOf, spy.Stdout.String()) {
				t.Errorf("expected stdout to be one of %#v, got %q", expectedStdoutOneOf, spy.Stdout.String())
			}
		})
	}
}
