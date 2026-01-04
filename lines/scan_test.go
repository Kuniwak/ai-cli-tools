package lines

import (
	"bufio"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestScanLinesWithNull(t *testing.T) {
	testCases := map[string]struct {
		input    string
		expected []string
	}{
		"empty": {
			input:    "",
			expected: []string{},
		},
		"one line": {
			input:    "one\u0000",
			expected: []string{"one"},
		},
		"several lines": {
			input:    "one\u0000two\u0000",
			expected: []string{"one", "two"},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(tc.input))
			scanner.Split(ScanLinesWithNull)
			actual := make([]string, 0)
			for scanner.Scan() {
				actual = append(actual, scanner.Text())
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Error(cmp.Diff(tc.expected, actual))
			}
		})
	}
}
