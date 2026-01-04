package split

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"

	"github.com/Kuniwak/ai-cli-tools/testableio"
)

type OutPathGenerator func(n int) string

func NewOutPathgenerator(outDir string, basenameTemplate string) OutPathGenerator {
	return func(n int) string {
		return filepath.Join(outDir, fmt.Sprintf(basenameTemplate, n))
	}
}

func SplitByLineCount(r io.Reader, lineCount int, outPathGenerator OutPathGenerator, openFileFunc testableio.OpenFileFunc) ([]string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	var n int
	writtenPaths := make([]string, 0)
	var w io.WriteCloser
	var err error
	for scanner.Scan() {
		if n%lineCount == 0 {
			if w != nil {
				if err := w.Close(); err != nil {
					return nil, fmt.Errorf("SplitByLineCount: failed to close previous writer: %w", err)
				}
			}
			path := outPathGenerator(n / lineCount)
			w, err = openFileFunc(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				return nil, fmt.Errorf("SplitByLineCount: failed to create new writer: %w", err)
			}
			writtenPaths = append(writtenPaths, path)
		}
		line := scanner.Text()
		fmt.Fprintln(w, line)
		n++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("SplitByLineCount: failed to scan lines: %w", err)
	}
	return writtenPaths, nil
}

func SplitByTotalCount(r io.Reader, totalCount int, outPathGenerator OutPathGenerator, openFileFunc testableio.OpenFileFunc) ([]string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	ls := make([]string, 0)
	for scanner.Scan() {
		ls = append(ls, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("SplitByTotalCount: failed to scan lines: %w", err)
	}
	lineCount := int(math.Ceil(float64(len(ls)) / float64(totalCount)))
	writtenPaths := make([]string, 0)
	for i := range totalCount {
		if err := func(i int) error {
			filePath := outPathGenerator(i)
			f, err := openFileFunc(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				return fmt.Errorf("SplitByTotalCount: failed to create new writer: %w", err)
			}
			defer f.Close()
			writtenPaths = append(writtenPaths, filePath)

			start := i * lineCount
			end := min(start+lineCount, len(ls))
			for _, line := range ls[start:end] {
				fmt.Fprintln(f, line)
			}

			return nil
		}(i); err != nil {
			return writtenPaths, err
		}
	}
	return writtenPaths, nil
}
