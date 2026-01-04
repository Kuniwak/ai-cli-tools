package testableio

import (
	"bytes"
	"io"
	"os"
)

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
