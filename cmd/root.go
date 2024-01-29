package cmd

import (
	"log"
	"os"
	"text/tabwriter"

	"github.com/mehix/k8s-resources/internal/aggr"
	"github.com/spf13/cobra"
)

var cmdRoot = &cobra.Command{
	Use:  "k8s-resources",
	Long: "Show a total of the resources declared in a set of yaml deployment files",
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, 0, 2, 3, ' ', tabwriter.AlignRight)
		defer w.Flush()

		if err := aggr.ShowAggregates(w, os.Stdin); err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() {
	if err := cmdRoot.Execute(); err != nil {
		log.Fatal(err)
	}
}
