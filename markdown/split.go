package markdown

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func NewWriterGenerator(outDir string, basenameTemplate string) func(n int) (io.WriteCloser, error) {
	return func(n int) (io.WriteCloser, error) {
		return os.OpenFile(filepath.Join(outDir, fmt.Sprintf(basenameTemplate, n)), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	}
}

type NopWriteCloser struct {
	w io.Writer
}

func NewNopWriteCloser(w io.Writer) *NopWriteCloser {
	return &NopWriteCloser{w: w}
}

func (n *NopWriteCloser) Write(p []byte) (int, error) {
	return n.w.Write(p)
}

func (n *NopWriteCloser) Close() error {
	return nil
}

func SpyWriterGenerator(bs *[]*bytes.Buffer) func(n int) (io.WriteCloser, error) {
	return func(n int) (io.WriteCloser, error) {
		*bs = append(*bs, &bytes.Buffer{})
		return NewNopWriteCloser((*bs)[len(*bs)-1]), nil
	}
}

// SplitBySections splits the reader into sections based on the number of seconds in each section.
// The writer generator function is called for each section to get the writer for the section.
func SplitBySections(r io.Reader, wGen func(n int) (io.WriteCloser, error)) error {
	var n int
	var w io.WriteCloser
	newWriterFunc := func() (io.WriteCloser, error) {
		if w != nil {
			if err := w.Close(); err != nil {
				return nil, fmt.Errorf("SplitBySections: failed to close previous writer: %w", err)
			}
		}
		w, err := wGen(n)
		if err != nil {
			return nil, fmt.Errorf("SplitBySections: failed to create new writer: %w", err)
		}
		n++
		return w, nil
	}

	w, err := newWriterFunc()
	if err != nil {
		return fmt.Errorf("SplitBySections: %w", err)
	}
	var hasPrevLine bool
	var prevLine string

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if hasPrevLine {
			if strings.HasPrefix(line, "===") || (strings.HasPrefix(line, "---") && prevLine != "") {
				w, err = newWriterFunc()
				if err != nil {
					return fmt.Errorf("SplitBySections: %w", err)
				}
				io.WriteString(w, prevLine)
				io.WriteString(w, "\n")
			} else if strings.HasPrefix(prevLine, "#") {
				w, err = newWriterFunc()
				if err != nil {
					return fmt.Errorf("SplitBySections: %w", err)
				}
				io.WriteString(w, prevLine)
				io.WriteString(w, "\n")
			} else {
				io.WriteString(w, prevLine)
				io.WriteString(w, "\n")
			}
		}
		hasPrevLine = true
		prevLine = line
	}
	if hasPrevLine {
		if strings.HasPrefix(prevLine, "#") {
			w, err = newWriterFunc()
			if err != nil {
				return fmt.Errorf("SplitBySections: %w", err)
			}
			io.WriteString(w, prevLine)
			io.WriteString(w, "\n")
		} else {
			io.WriteString(w, prevLine)
			io.WriteString(w, "\n")
		}
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("SplitBySections: %w", err)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("SplitBySections: %w", err)
	}
	return nil
}
