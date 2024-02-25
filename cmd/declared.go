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

var (
	dataSrc    string
	output     string
	onlyTotals bool
)

var cmdDeclared = &cobra.Command{
	Use:   "declared [--src {confluentinc | k8s}]",
	Long:  "Read yaml files from standard input and print out resources declared and their total values",
	Short: "Show a total of the resources declared in a set of yaml deployment files",
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, 0, 2, 3, ' ', tabwriter.AlignRight)
		defer w.Flush()

		switch dataSrc {
		case "confluentinc":
			aggregator := aggr.New(confluentinc.ComputeAggregates)
			if err := aggregator.Load(os.Stdin); err != nil {
				log.Fatal(err)
			}
			aggregator.Print(w, []string{"Kind", "repl", "cpuR", "memR", "cpuL", "memL", "vol"}, onlyTotals)
		case "k8s":
			aggregator := aggr.New(k8s.ComputeAggregates)
			if err := aggregator.Load(os.Stdin); err != nil {
				log.Fatal(err)
			}
			aggregator.Print(w, []string{"Kind", "repl", "cpuR", "memR", "cpuL", "memL"}, onlyTotals)
		default:
			fmt.Println("Unknown data source: ", dataSrc)
			os.Exit(1)
		}
	},
}

func init() {
	cmdDeclared.PersistentFlags().StringVar(&dataSrc, "src", "k8s", "Source of the YAML files")
	cmdDeclared.Flags().StringVarP(&output, "output", "o", "plain", "Output format (plain, json)")
	cmdDeclared.Flags().BoolVar(&onlyTotals, "onlyTotals", false, "Print only the totals")
}
