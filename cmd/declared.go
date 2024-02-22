package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/mehix/kuberes/internal/aggr"
	"github.com/mehix/kuberes/internal/aggr/confluentinc"
	"github.com/mehix/kuberes/internal/aggr/k8s"
	"github.com/spf13/cobra"
)

var dataSrc string

var cmdDeclared = &cobra.Command{
	Use:   "declared [--src {confluentinc | k8s}]",
	Long:  "Read yaml files from standard input and print out resources declared and their total values",
	Short: "Show a total of the resources declared in a set of yaml deployment files",
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, 0, 2, 3, ' ', tabwriter.AlignRight)
		defer w.Flush()

		switch dataSrc {
		case "confluentinc":
			headers := []string{"Kind", "repl", "cpuR", "memR", "cpuL", "memL", "vol"}
			aggregator := aggr.NewAggregator[confluentinc.Objects](headers, confluentinc.ComputeAggregates)
			if err := aggregator.PrintResources(w, os.Stdin); err != nil {
				log.Fatal(err)
			}
		case "k8s":
			headers := []string{"Name", "Kind", "repl", "cpuR", "memR", "cpuL", "memL"}
			aggregator := aggr.NewAggregator[k8s.Objects](headers, k8s.ComputeAggregates)
			if err := aggregator.PrintResources(w, os.Stdin); err != nil {
				log.Fatal(err)
			}
		default:
			fmt.Println("Unknown data source: ", dataSrc)
			os.Exit(1)
		}
	},
}

func init() {
	cmdDeclared.PersistentFlags().StringVar(&dataSrc, "src", "k8s", "Source of the YAML files")
}
