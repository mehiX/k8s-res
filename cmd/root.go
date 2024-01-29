package cmd

import (
	"log"
	"os"

	"github.com/mehix/k8s-resources/pkg/aggr"
	"github.com/spf13/cobra"
)

var cmdRoot = &cobra.Command{
	Use:  "k8s-resources",
	Long: "Show a total of the resources declared in a set of yaml deployment files",
	Run: func(cmd *cobra.Command, args []string) {
		if err := aggr.ShowAggregates(os.Stdin); err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() {
	if err := cmdRoot.Execute(); err != nil {
		log.Fatal(err)
	}
}
