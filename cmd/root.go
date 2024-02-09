package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	version, commit, date = "dev", "dev", "n/a"
)

var cmdRoot = &cobra.Command{
	Use:   "kres show",
	Short: "Show sizes of Kubernetes resources",
}

func Execute() {
	cmdRoot.AddCommand(cmdVersion, cmdDeclared, cmdNamespaces)

	if err := cmdRoot.Execute(); err != nil {
		log.Fatal(err)
	}
}
