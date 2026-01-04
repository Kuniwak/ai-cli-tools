package markdown

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Kuniwak/ai-cli-tools/split"
	"github.com/Kuniwak/ai-cli-tools/testableio"
)

// SplitBySections splits the reader into sections based on the number of seconds in each section.
// The writer generator function is called for each section to get the writer for the section.
func SplitBySections(r io.Reader, outPathGenerator split.OutPathGenerator, openFileFunc testableio.OpenFileFunc) ([]string, error) {
	var n int
	var w io.WriteCloser
	writtenPaths := make([]string, 0)
	newWriterFunc := func() (io.WriteCloser, error) {
		if w != nil {
			if err := w.Close(); err != nil {
				return nil, fmt.Errorf("SplitBySections: failed to close previous writer: %w", err)
			}
		}
		path := outPathGenerator(n)
		w, err := openFileFunc(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return nil, fmt.Errorf("SplitBySections: failed to create new writer: %w", err)
		}
		writtenPaths = append(writtenPaths, path)
		n++
		return w, nil
	}

	w, err := newWriterFunc()
	if err != nil {
		return writtenPaths, fmt.Errorf("SplitBySections: %w", err)
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
					return writtenPaths, fmt.Errorf("SplitBySections: %w", err)
				}
				io.WriteString(w, prevLine)
				io.WriteString(w, "\n")
			} else if strings.HasPrefix(prevLine, "#") {
				w, err = newWriterFunc()
				if err != nil {
					return writtenPaths, fmt.Errorf("SplitBySections: %w", err)
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
				return writtenPaths, fmt.Errorf("SplitBySections: %w", err)
			}
			io.WriteString(w, prevLine)
			io.WriteString(w, "\n")
		} else {
			io.WriteString(w, prevLine)
			io.WriteString(w, "\n")
		}
	}
	if err := w.Close(); err != nil {
		return writtenPaths, fmt.Errorf("SplitBySections: %w", err)
	}
	if err := scanner.Err(); err != nil {
		return writtenPaths, fmt.Errorf("SplitBySections: %w", err)
	}
	return writtenPaths, nil
}
