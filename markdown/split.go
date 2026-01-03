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

type OutPathGenerator func(n int) string

func NewOutPathgenerator(outDir string, basenameTemplate string) OutPathGenerator {
	return func(n int) string {
		return filepath.Join(outDir, fmt.Sprintf(basenameTemplate, n))
	}
}

type OpenFileFunc func(path string, flag int, perm os.FileMode) (io.WriteCloser, error)

func NewOpenFileFunc() OpenFileFunc {
	return func(path string, flag int, perm os.FileMode) (io.WriteCloser, error) {
		return os.OpenFile(path, flag, perm)
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

type SpyOpenFileFunc struct {
	m map[string]*bytes.Buffer
}

func NewSpyOpenFileFunc() *SpyOpenFileFunc {
	return &SpyOpenFileFunc{m: make(map[string]*bytes.Buffer)}
}

func (s *SpyOpenFileFunc) Written() map[string]string {
	m := make(map[string]string)
	for path, w := range s.m {
		m[path] = w.String()
	}
	return m
}

func (s *SpyOpenFileFunc) OpenFileFunc() OpenFileFunc {
	return func(path string, _ int, _ os.FileMode) (io.WriteCloser, error) {
		w, ok := s.m[path]
		if !ok {
			w = bytes.NewBuffer(nil)
			s.m[path] = w
		}
		return NewNopWriteCloser(w), nil
	}
}

// SplitBySections splits the reader into sections based on the number of seconds in each section.
// The writer generator function is called for each section to get the writer for the section.
func SplitBySections(r io.Reader, outPathGenerator OutPathGenerator, openFileFunc OpenFileFunc) ([]string, error) {
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
