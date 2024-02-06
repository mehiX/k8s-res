package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var dataSrc string

var cmdRoot = &cobra.Command{
	Use:   "k8s-res",
	Short: "Show sizes of Kubernetes resources",
}

func Execute() {
	cmdRoot.AddCommand(cmdShow, cmdNamespaces)

	if err := cmdRoot.Execute(); err != nil {
		log.Fatal(err)
	}
}
