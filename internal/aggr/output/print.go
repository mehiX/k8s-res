package output

import (
	"fmt"
	"io"
	"strings"
)

type Printable interface {
	Outputline() string
	IsEmpty() bool
}

func toTxt(objs []Printable) (lines []string) {

	for _, o := range objs {

		if o.IsEmpty() {
			continue
		}

		lines = append(lines, o.Outputline())
	}

	return
}

func Print(w io.Writer, objs []Printable, headers []string) {

	delimiter := strings.Repeat("--------\t", len(headers))

	lines := toTxt(objs)

	fmt.Fprintln(w, strings.Join(headers, "\t")+"\t")
	fmt.Fprintln(w, delimiter)

	for i, l := range lines {
		// last element consists of aggregated values
		if i == len(lines)-1 {
			fmt.Fprintln(w, delimiter)
		}
		fmt.Fprintln(w, l)
	}
}
