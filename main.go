/*
Read YAML input (ex: the output of `helm template`) and output the resources needed
based on the resources declarations and replicas
*/
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/api/resource"
)

func main() {

	var buf bytes.Buffer
	_, err := io.Copy(&buf, os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	parts := bytes.Split(buf.Bytes(), []byte("---"))

	w := tabwriter.NewWriter(os.Stdout, 0, 2, 3, ' ', tabwriter.AlignRight)
	defer w.Flush()
	fmt.Fprintln(w, "Kind\trepl\tcpuR\tmemR\tcpuL\tmemL\t")

	delimiter := "--------\t-------\t-------\t-------\t------\t------\t"
	fmt.Fprintln(w, delimiter)

	tCpuR := new(resource.Quantity)
	tCpuL := new(resource.Quantity)
	tMemR := new(resource.Quantity)
	tMemL := new(resource.Quantity)
	for _, p := range parts {
		var o Object
		if err := yaml.NewDecoder(bytes.NewReader(p)).Decode(&o); err != nil && err != io.EOF {
			log.Println(err)
		} else {
			// missing replicas declaration, we assume 1 replica
			if o.Spec.Replicas == 0 {
				o.Spec.Replicas = 1
			}

			var cpuR resource.Quantity
			if v := o.Spec.PodTemplate.Resources.Requests["cpu"]; v != "" {
				cpuR = resource.MustParse(v)
			}

			var memR resource.Quantity
			if v := o.Spec.PodTemplate.Resources.Requests["memory"]; v != "" {
				memR = resource.MustParse(v)
			}

			var cpuL resource.Quantity
			if v := o.Spec.PodTemplate.Resources.Limits["cpu"]; v != "" {
				cpuL = resource.MustParse(v)
			}

			var memL resource.Quantity
			if v := o.Spec.PodTemplate.Resources.Limits["memory"]; v != "" {
				memL = resource.MustParse(v)
			}

			for i := 0; i < o.Spec.Replicas; i++ {
				tCpuR.Add(cpuR)
				tCpuL.Add(cpuL)
				tMemR.Add(memR)
				tMemL.Add(memL)
			}
			fmt.Fprintf(
				w,
				"%s\t%d\t%s\t%s\t%s\t%s\t\n",
				o.Kind, o.Spec.Replicas, cpuR.String(), memR.String(), cpuL.String(), memL.String())
		}

	}

	fmt.Fprintln(w, delimiter)
	fmt.Fprintf(
		w,
		"\t\t%s\t%s\t%s\t%s\t\n",
		tCpuR.String(), tMemR.String(), tCpuL.String(), tMemL.String())
}

type Object struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Spec       spec   `yaml:"spec"`
}

type spec struct {
	PodTemplate podTemplate `yaml:"podTemplate"`
	Replicas    int         `yaml:"replicas"`
}

type podTemplate struct {
	Resources resources `yaml:"resources"`
}

type resources struct {
	Limits   map[string]string `yaml:"limits"`
	Requests map[string]string `yaml:"requests"`
}
