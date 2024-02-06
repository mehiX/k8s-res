package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/mehix/k8s-res/internal/namespaces"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var kCtx *string
var kconf *string
var onlyNamespaces *string
var onlyNamespacesFromFile *string

func init() {

	if home := homedir.HomeDir(); home != "" {
		kconf = cmdNamespaces.PersistentFlags().String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kconf = cmdNamespaces.PersistentFlags().String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	kCtx = cmdNamespaces.PersistentFlags().String("context", "", "Choose a Kubernetes context. If left empty, then the default will be used")

	onlyNamespaces = cmdNamespaces.PersistentFlags().String("only", "", "comma-separated list of namespaces")

	onlyNamespacesFromFile = cmdNamespaces.PersistentFlags().String("only-file", "", "Location of file with the list of namespaces to query")
}

var cmdNamespaces = &cobra.Command{
	Use:   "namespaces",
	Short: "Show resources for namespace(s)",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := Resources(ctx); err != nil {
			log.Fatal(err)
		}
	},
}

func whichNamespaces() ([]string, error) {

	lst := *onlyNamespaces

	if *onlyNamespacesFromFile != "" {
		b, err := os.ReadFile(*onlyNamespacesFromFile)
		if err != nil {
			return nil, fmt.Errorf("Reading input file %s got error: %w", *onlyNamespacesFromFile, err)
		}
		b = bytes.ReplaceAll(b, []byte("\n"), []byte(","))
		lst += "," + string(b)
	}

	ns := strings.Split(lst, ",")
	for i := range ns {
		ns[i] = strings.TrimSpace(ns[i])
	}
	slices.Sort(ns)
	ns = slices.Compact(ns)

	if ns[0] == "" {
		// if there is an empty element then it is in the first position
		ns = ns[1:]
	}

	return ns, nil
}

func Resources(ctx context.Context) error {

	filterByNS, err := whichNamespaces()
	if err != nil {
		return err
	}

	configLoadRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: *kconf}
	configOverrides := &clientcmd.ConfigOverrides{CurrentContext: *kCtx}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadRules, configOverrides).ClientConfig()
	if err != nil {
		return err
	}

	u, err := url.Parse(config.Host)
	if err == nil {
		fmt.Println("Context: ", strings.Split(u.Host, ".")[0])
		fmt.Println()
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	quotas, err := namespaces.GetQuotas(ctx, clientset, filterByNS)
	if err != nil {
		return err
	}

	printQuotas(os.Stdout, quotas)

	return nil
}

func printQuotas(out io.Writer, results []namespaces.QuotaResult) {

	if len(results) == 0 {
		fmt.Println("no data")
		return
	}

	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', tabwriter.AlignRight)
	defer w.Flush()

	delim := strings.Repeat("---------\t", 6)
	lf := "%s\t%v\t%v\t%v\t%v\t%v\t\n"

	fmt.Fprintln(w, "namespace\tcpuR\tcpuL\tmemoryR\tmemoryL\tstorage\t")
	fmt.Fprintln(w, delim)

	var cpuRT, cpuLT, memRT, memLT, storageT = new(resource.Quantity),
		new(resource.Quantity),
		new(resource.Quantity),
		new(resource.Quantity),
		new(resource.Quantity)

	for _, r := range results {
		if r.Error != nil {
			log.Println(r.Error.Error())
			continue
		}

		for _, q := range r.Quotas {
			h := q.Spec.Hard

			cpuRT.Add(*h.Name(v1.ResourceRequestsCPU, resource.DecimalSI))
			cpuLT.Add(*h.Name(v1.ResourceLimitsCPU, resource.DecimalSI))
			memRT.Add(*h.Name(v1.ResourceRequestsMemory, resource.BinarySI))
			memLT.Add(*h.Name(v1.ResourceLimitsMemory, resource.BinarySI))
			storageT.Add(*h.Name(v1.ResourceRequestsStorage, resource.DecimalSI))

			fmt.Fprintf(w, lf,
				q.Namespace,
				h.Name(v1.ResourceRequestsCPU, resource.DecimalSI),
				h.Name(v1.ResourceLimitsCPU, resource.DecimalSI),
				h.Name(v1.ResourceRequestsMemory, resource.BinarySI),
				h.Name(v1.ResourceLimitsMemory, resource.BinarySI),
				h.Name(v1.ResourceRequestsStorage, resource.DecimalSI),
			)
		}
	}

	fmt.Fprintln(w, delim)
	fmt.Fprintf(w, lf, "total",
		cpuRT,
		cpuLT,
		memRT,
		memLT,
		storageT)

	fmt.Fprintln(w)
}
