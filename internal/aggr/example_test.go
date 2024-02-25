package aggr

import (
	"bytes"
	_ "embed"
	"log"
	"os"
	"text/tabwriter"

	"github.com/mehix/kuberes/internal/aggr/confluentinc"
	"github.com/mehix/kuberes/internal/aggr/k8s"
)

//go:embed testdata_k8s.yaml
var dataK8s []byte

//go:embed testdata_confluentinc.yaml
var dataConfluentinc []byte

func ExamplePrinter_PrintResources_k8s() {

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, '.', tabwriter.AlignRight|tabwriter.Debug)
	defer w.Flush()

	a := New(k8s.ComputeAggregates)
	if err := a.Load(bytes.NewReader(dataK8s)); err != nil {
		log.Fatal(err)
	}
	a.Print(w, []string{"Kind", "repl", "cpuR", "memR", "cpuL", "memL"}, false)

	// output:
	// .............................Name|......Kind|....repl|....cpuR|....memR|....cpuL|....memL|
	// .........................--------|..--------|--------|--------|--------|--------|--------|
	// ................release-name-test|Deployment|.......4|....100m|...128Mi|....200m|...256Mi|
	// release-name-test-test-connection|.......Pod|.......2|.....50m|...100Mi|....100m|...150Mi|
	// .........................--------|..--------|--------|--------|--------|--------|--------|
	// .................................|..........|........|....500m|...712Mi|.......1|..1324Mi|
}

func ExamplePrinter_PrintResources_k8s_onlyTotals() {

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, '.', tabwriter.AlignRight|tabwriter.Debug)
	defer w.Flush()

	a := New(k8s.ComputeAggregates)
	if err := a.Load(bytes.NewReader(dataK8s)); err != nil {
		log.Fatal(err)
	}
	a.Print(w, []string{"Kind", "repl", "cpuR", "memR", "cpuL", "memL"}, true)

	// output:
	// ....Kind|....repl|....cpuR|....memR|....cpuL|....memL|
	// --------|--------|--------|--------|--------|--------|
	// ........|........|....500m|...712Mi|.......1|..1324Mi|
}

func ExamplePrinter_PrintResources_confluentinc() {

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, '.', tabwriter.AlignRight|tabwriter.Debug)
	defer w.Flush()

	a := New(confluentinc.ComputeAggregates)
	if err := a.Load(bytes.NewReader(dataConfluentinc)); err != nil {
		log.Fatal(err)
	}
	a.Print(w, []string{"Kind", "repl", "cpuR", "memR", "cpuL", "memL", "vol"}, false)

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

func ExamplePrinter_PrintResources_confluentinc_onlyTotals() {

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, '.', tabwriter.AlignRight|tabwriter.Debug)
	defer w.Flush()

	a := New(confluentinc.ComputeAggregates)
	if err := a.Load(bytes.NewReader(dataConfluentinc)); err != nil {
		log.Fatal(err)
	}
	a.Print(w, []string{"Kind", "repl", "cpuR", "memR", "cpuL", "memL", "vol"}, true)

	// Output:
	// ....Kind|....repl|....cpuR|....memR|....cpuL|....memL|.....vol|
	// --------|--------|--------|--------|--------|--------|--------|
	// ........|........|...2750m|.41936Mi|...3400m|.48056Mi|...450Gi|
}
