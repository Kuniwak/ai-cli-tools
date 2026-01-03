package markdown

import (
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSplitBySections(t *testing.T) {
	testCases := map[string]struct {
		input    string
		expected map[string]string
	}{
		"empty": {
			input: "",
			expected: map[string]string{
				"out/test-0.md": "",
			},
		},
		"horizontal bar": {
			input: `aaa

---

bbb
`,
			expected: map[string]string{
				"out/test-0.md": `aaa

---

bbb
`,
			},
		},
		"section with hashes": {
			input: `# 1`,
			expected: map[string]string{
				"out/test-0.md": "",
				"out/test-1.md": `# 1
`,
			},
		},
		"section with equals": {
			input: `aaa
===
`,
			expected: map[string]string{
				"out/test-0.md": "",
				"out/test-1.md": `aaa
===
`,
			},
		},
		"section with dashes": {
			input: `aaa
---
`,
			expected: map[string]string{
				"out/test-0.md": "",
				"out/test-1.md": `aaa
---
`},
		},
		"section with hashes and content before the section": {
			input: `aaa
# bbb`,
			expected: map[string]string{
				"out/test-0.md": `aaa
`,
				"out/test-1.md": `# bbb
`},
		},
		"section with equals and content before the section": {
			input: `aaa
bbb
===
`,
			expected: map[string]string{
				"out/test-0.md": `aaa
`,
				"out/test-1.md": `bbb
===
`},
		},
		"section with dashes and content before the section": {
			input: `aaa
bbb
---
`,
			expected: map[string]string{
				"out/test-0.md": `aaa
`,
				"out/test-1.md": `bbb
---
`},
		},
		"section with hashes and content after the section": {
			input: `# aaa
bbb`,
			expected: map[string]string{
				"out/test-0.md": "",
				"out/test-1.md": `# aaa
bbb
`},
		},
		"section with equals and content after the section": {
			input: `aaa
===
bbb`,
			expected: map[string]string{
				"out/test-0.md": "",
				"out/test-1.md": `aaa
===
bbb
`},
		},
		"section with dashes and content after the section": {
			input: `aaa
---
bbb`,
			expected: map[string]string{
				"out/test-0.md": "",
				"out/test-1.md": `aaa
---
bbb
`},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			outPathGenerator := NewOutPathgenerator("out", "test-%d.md")
			spy := NewSpyOpenFileFunc()
			_, err := SplitBySections(strings.NewReader(tc.input), outPathGenerator, spy.OpenFileFunc())
			if err != nil {
				t.Fatalf("SplitBySections: %v", err)
			}
			written := spy.Written()
			if !reflect.DeepEqual(written, tc.expected) {
				t.Error(cmp.Diff(tc.expected, written))
			}
		})
	}
}
