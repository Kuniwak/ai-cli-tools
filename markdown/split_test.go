package markdown

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSplitBySections(t *testing.T) {
	testCases := map[string]struct {
		input    string
		expected []string
	}{
		"empty": {
			input:    "",
			expected: []string{""},
		},
		"horizontal bar": {
			input: `aaa

---

bbb
`,
			expected: []string{
				`aaa

---

bbb
`},
		},
		"section with hashes": {
			input: `# 1`,
			expected: []string{
				"",
				`# 1
`},
		},
		"section with equals": {
			input: `aaa
===
`,
			expected: []string{
				"",
				`aaa
===
`},
		},
		"section with dashes": {
			input: `aaa
---
`,
			expected: []string{
				"",
				`aaa
---
`},
		},
		"section with hashes and content before the section": {
			input: `aaa
# bbb`,
			expected: []string{
				`aaa
`,
				`# bbb
`},
		},
		"section with equals and content before the section": {
			input: `aaa
bbb
===
`,
			expected: []string{
				`aaa
`,
				`bbb
===
`},
		},
		"section with dashes and content before the section": {
			input: `aaa
bbb
---
`,
			expected: []string{
				`aaa
`,
				`bbb
---
`},
		},
		"section with hashes and content after the section": {
			input: `# aaa
bbb`,
			expected: []string{
				"",
				`# aaa
bbb
`},
		},
		"section with equals and content after the section": {
			input: `aaa
===
bbb`,
			expected: []string{
				"",
				`aaa
===
bbb
`},
		},
		"section with dashes and content after the section": {
			input: `aaa
---
bbb`,
			expected: []string{
				"",
				`aaa
---
bbb
`},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			bs := []*bytes.Buffer{}
			wGen := SpyWriterGenerator(&bs)
			err := SplitBySections(strings.NewReader(tc.input), wGen)
			if err != nil {
				t.Fatalf("SplitBySections: %v", err)
			}
			bs2 := make([]string, len(bs))
			for i, b := range bs {
				bs2[i] = b.String()
			}
			if !reflect.DeepEqual(bs2, tc.expected) {
				t.Error(cmp.Diff(tc.expected, bs2))
			}
		})
	}
}
