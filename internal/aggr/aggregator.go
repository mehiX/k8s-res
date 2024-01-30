package aggr

import (
	"bytes"
	"io"
	"log"

	"gopkg.in/yaml.v3"
)

type Object = Printable

type Aggregator[T ~[]O, O Object] struct {
	Headers []string
	f       func(T) T // aggregator function
}

func NewAggregator[T ~[]O, O Object](headers []string, aggrFunc func(T) T) Aggregator[T, O] {
	return Aggregator[T, O]{
		Headers: headers,
		f:       aggrFunc,
	}
}

func (a Aggregator[T, O]) PrintResources(w io.Writer, in io.Reader) error {

	lines, err := a.getWithAggregates(in)
	if err != nil {
		return err
	}

	printAll(w, lines, a.Headers)

	return nil
}

func (a Aggregator[T, O]) getWithAggregates(in io.Reader) ([]Printable, error) {

	allObjects, err := a.readObjects(in)
	if err != nil {
		return nil, err
	}

	allObjects = a.f(allObjects)

	printItems := make([]Printable, len(allObjects))
	for i := range allObjects {
		printItems[i] = Printable(allObjects[i])
	}

	return printItems, nil
}

func (a Aggregator[T, O]) readObjects(in io.Reader) (T, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, in)
	if err != nil {
		return nil, err
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

	return allObjects, nil
}
