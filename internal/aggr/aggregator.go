package aggr

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"

	"gopkg.in/yaml.v3"
)

type Aggregator[T ~[]O, O Printable] struct {
	objects       T
	aggregated    O
	aggregateFunc func(T) O
}

func New[T ~[]O, O Printable](aggrFunc func(T) O) *Aggregator[T, O] {
	return &Aggregator[T, O]{
		aggregateFunc: aggrFunc,
	}
}

func (a *Aggregator[T, O]) Load(in io.Reader) error {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, in)
	if err != nil {
		return err
	}

	parts := bytes.Split(buf.Bytes(), []byte("---"))

	var allObjects T
	for _, p := range parts {
		var o O
		if err := yaml.NewDecoder(bytes.NewReader(p)).Decode(&o); err != nil && err != io.EOF {
			log.Println(err)
		} else {
			allObjects = append(allObjects, o)
		}
	}

	a.objects = allObjects
	a.aggregated = a.aggregateFunc(allObjects)

	return nil
}

func (a *Aggregator[T, O]) Print(w io.Writer, headers []string, onlyTotals bool) {

	delimiter := strings.Repeat("--------\t", len(headers))

	fmt.Fprintln(w, strings.Join(headers, "\t")+"\t")
	fmt.Fprintln(w, delimiter)

	if !onlyTotals {
		for _, l := range a.objects {
			if !l.IsEmpty() {
				fmt.Fprintln(w, l.String())
			}
		}
		fmt.Fprintln(w, delimiter)
	}

	fmt.Fprintln(w, a.aggregated.String())

	fmt.Fprintln(w)
}
