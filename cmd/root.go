package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/mehix/k8s-resources/internal/aggr"
	"github.com/mehix/k8s-resources/internal/aggr/confluentinc"
	"github.com/mehix/k8s-resources/internal/aggr/k8s"
	"github.com/spf13/cobra"
)

var dataSrc string

var cmdRoot = &cobra.Command{
	Use:  "k8s-resources",
	Long: "Show a total of the resources declared in a set of yaml deployment files",
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, 0, 2, 3, ' ', tabwriter.AlignRight)
		defer w.Flush()

		switch dataSrc {
		case "confluentinc":
			headers := []string{"Kind", "repl", "cpuR", "memR", "cpuL", "memL", "vol"}
			aggregator := aggr.NewAggregator[confluentinc.Objects](headers, confluentinc.AddAggregates)
			if err := aggregator.PrintResources(w, os.Stdin); err != nil {
				log.Fatal(err)
			}
		case "k8s":
			headers := []string{"Kind", "repl", "cpuR", "memR", "cpuL", "memL"}
			aggregator := aggr.NewAggregator[k8s.Objects](headers, k8s.AddAggregates)
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
	cmdRoot.PersistentFlags().StringVar(&dataSrc, "src", "confluentinc", "Source of the YAML files")
}

func Execute() {
	if err := cmdRoot.Execute(); err != nil {
		log.Fatal(err)
	}
}
