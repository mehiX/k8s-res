package aggr

import (
	"bytes"
	_ "embed"
	"os"
	"text/tabwriter"

	"github.com/mehix/k8s-resources/internal/aggr/confluentinc"
	"github.com/mehix/k8s-resources/internal/aggr/k8s"
)

//go:embed testdata_k8s.yaml
var dataK8s []byte

//go:embed testdata_confluentinc.yaml
var dataConfluentinc []byte

func ExampleAggregator_PrintResources_k8s() {

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, '.', tabwriter.AlignRight|tabwriter.Debug)
	defer w.Flush()

	headers := []string{"Kind", "repl", "cpuR", "memR", "cpuL", "memL"}
	aggregator := NewAggregator[k8s.Objects](headers, k8s.ComputeAggregates)

	aggregator.PrintResources(w, bytes.NewReader(dataK8s))
	// output:
	// ......Kind|....repl|....cpuR|....memR|....cpuL|....memL|
	// ..--------|--------|--------|--------|--------|--------|
	// Deployment|.......4|....100m|...128Mi|....200m|...256Mi|
	// .......Pod|.......2|.....50m|...100Mi|....100m|...150Mi|
	// ..--------|--------|--------|--------|--------|--------|
	// ..........|........|....500m|...712Mi|.......1|..1324Mi|
}

func ExampleAggregator_PrintResources_confluentinc() {

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, '.', tabwriter.AlignRight|tabwriter.Debug)
	defer w.Flush()

	headers := []string{"Kind", "repl", "cpuR", "memR", "cpuL", "memL", "vol"}
	aggregator := NewAggregator[confluentinc.Objects](headers, confluentinc.ComputeAggregates)

	aggregator.PrintResources(w, bytes.NewReader(dataConfluentinc))
	// Output:
	// ..........Kind|....repl|....cpuR|....memR|....cpuL|....memL|.....vol|
	// ......--------|--------|--------|--------|--------|--------|--------|
	// .......Connect|.......1|....400m|.....5Gi|....450m|.....6Gi|.......0|
	// .ControlCenter|.......1|...1100m|....12Gi|...1200m|....13Gi|...150Gi|
	// .........Kafka|.......3|....250m|.....4Gi|....300m|.....5Gi|....50Gi|
	// KafkaRestProxy|.......2|.....50m|..1000Mi|....100m|..1500Mi|.......0|
	// SchemaRegistry|.......2|.....50m|.....2Gi|....100m|.....2Gi|.......0|
	// .....Zookeeper|.......3|....100m|.....2Gi|....150m|.....2Gi|....50Gi|
	// ......--------|--------|--------|--------|--------|--------|--------|
	// ..............|........|...2750m|.41936Mi|...3400m|.48056Mi|...450Gi|
}
