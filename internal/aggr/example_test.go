package aggr

import (
	"bytes"
	_ "embed"
	"os"
	"text/tabwriter"
)

//go:embed testdata.yaml
var data []byte

func ExampleShowAggregates() {

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, '.', tabwriter.AlignRight|tabwriter.Debug)
	defer w.Flush()

	ShowAggregates(w, bytes.NewReader(data))
	// Output:
	// ..........Kind|...repl|...cpuR|...memR|..cpuL|...memL|....vol|
	// ......--------|-------|-------|-------|------|.------|-------|
	// .......Connect|......1|...400m|....5Gi|..450m|....6Gi|.......|
	// .ControlCenter|......1|..1100m|...12Gi|.1200m|...13Gi|..150Gi|
	// .........Kafka|......3|...250m|....4Gi|..300m|....5Gi|...50Gi|
	// KafkaRestProxy|......2|....50m|.1000Mi|..100m|.1500Mi|.......|
	// SchemaRegistry|......2|....50m|....2Gi|..100m|....2Gi|.......|
	// .....Zookeeper|......3|...100m|....2Gi|..150m|....2Gi|...40Gi|
	// ......--------|-------|-------|-------|------|.------|-------|
	// ..............|......0|..2750m|41936Mi|.3400m|48056Mi|..420Gi|
}
