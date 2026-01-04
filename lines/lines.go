package lines

import (
	"fmt"
	"io"
)

func WriteLine(null bool, line string, w io.Writer) error {
	if null {
		io.WriteString(w, line)
		io.WriteString(w, "\000")
	} else {
		fmt.Fprintln(w, line)
	}
	return nil
}

func WriteLines(null bool, lines []string, w io.Writer) error {
	for _, line := range lines {
		if err := WriteLine(null, line, w); err != nil {
			return err
		}
	}
	return nil
}
