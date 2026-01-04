package split

import (
	"strings"
	"testing"

	"github.com/Kuniwak/ai-cli-tools/testableio"
)

func TestSplitByLineCount(t *testing.T) {
	testCases := map[string]struct {
		input     string
		lineCount int
		expected  map[string]string
	}{
		"empty": {
			input:     "",
			lineCount: 1,
			expected:  map[string]string{"out/test-0.txt": ""},
		},
		"write to one file": {
			input:     "one\n",
			lineCount: 1,
			expected:  map[string]string{"out/test-0.txt": "one\n"},
		},
		"write to several files": {
			input:     "one\ntwo\n",
			lineCount: 1,
			expected:  map[string]string{"out/test-0.txt": "one\n", "out/test-1.txt": "two\n"},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			outPathGenerator := NewOutPathgenerator("out", "test-%d.txt")
			spy := testableio.NewSpyOpenFileFunc()
			_, err := SplitByLineCount(strings.NewReader(tc.input), tc.lineCount, outPathGenerator, spy.OpenFileFunc())
			if err != nil {
				t.Fatalf("SplitByLineCount: %v", err)
			}
		})
	}
}

func TestSplitByTotalCount(t *testing.T) {
	testCases := map[string]struct {
		input      string
		totalCount int
		expected   map[string]string
	}{
		"empty": {
			input:      "",
			totalCount: 2,
			expected:   map[string]string{"out/test-0.txt": "", "out/test-1.txt": ""},
		},
		"write to one file": {
			input:      "one\n",
			totalCount: 2,
			expected:   map[string]string{"out/test-0.txt": "one\n", "out/test-1.txt": ""},
		},
		"write to several files": {
			input:      "one\ntwo\n",
			totalCount: 2,
			expected:   map[string]string{"out/test-0.txt": "one\n", "out/test-1.txt": "two\n"},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			outPathGenerator := NewOutPathgenerator("out", "test-%d.txt")
			spy := testableio.NewSpyOpenFileFunc()
			_, err := SplitByTotalCount(strings.NewReader(tc.input), tc.totalCount, outPathGenerator, spy.OpenFileFunc())
			if err != nil {
				t.Fatalf("SplitByTotalCount: %v", err)
			}
		})
	}
}
