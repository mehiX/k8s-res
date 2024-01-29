package aggr

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

func ShowAggregates(in io.Reader) error {
	allObjects, err := readObjects(in)
	if err != nil {
		return err
	}

	lines := toOuputWithAggregates(allObjects)

	print(lines)

	return nil
}

func readObjects(in io.Reader) ([]Object, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, in)
	if err != nil {
		return nil, err
	}

	parts := bytes.Split(buf.Bytes(), []byte("---"))

	var allObjects []Object
	for _, p := range parts {
		var o Object
		if err := yaml.NewDecoder(bytes.NewReader(p)).Decode(&o); err != nil && err != io.EOF {
			log.Println(err)
		} else {
			allObjects = append(allObjects, o)
		}
	}

	return allObjects, nil
}

func print(lines []OutputLine) {

	w := tabwriter.NewWriter(os.Stdout, 0, 2, 3, ' ', tabwriter.AlignRight)
	defer w.Flush()
	fmt.Fprintln(w, "Kind\trepl\tcpuR\tmemR\tcpuL\tmemL\tvol\t")

	delimiter := "--------\t-------\t-------\t-------\t------\t------\t-------\t"
	fmt.Fprintln(w, delimiter)

	for i, l := range lines {
		// last element consists of aggregated values
		if i == len(lines)-1 {
			fmt.Fprintln(w, delimiter)
		}
		fmt.Fprintln(w, l.String())
	}
}

func calculateAggregates(objs []Object) Object {

	tCpuR := new(resource.Quantity)
	tCpuL := new(resource.Quantity)
	tMemR := new(resource.Quantity)
	tMemL := new(resource.Quantity)
	dataVol := new(resource.Quantity)

	for _, o := range objs {
		if o.Spec.PodTemplate == nil {
			continue
		}

		// missing replicas declaration, we assume 1 replica
		if o.Spec.Replicas == 0 {
			o.Spec.Replicas = 1
		}

		var cpuR resource.Quantity
		if v := o.Spec.PodTemplate.Resources.Requests.Cpu; v != "" {
			cpuR = resource.MustParse(v)
		}

		var memR resource.Quantity
		if v := o.Spec.PodTemplate.Resources.Requests.Memory; v != "" {
			memR = resource.MustParse(v)
		}

		var cpuL resource.Quantity
		if v := o.Spec.PodTemplate.Resources.Limits.Cpu; v != "" {
			cpuL = resource.MustParse(v)
		}

		var memL resource.Quantity
		if v := o.Spec.PodTemplate.Resources.Limits.Memory; v != "" {
			memL = resource.MustParse(v)
		}

		var dv resource.Quantity
		if v := o.Spec.DataVolumeCapacity; v != "" {
			dv = resource.MustParse(v)
		}

		for i := 0; i < o.Spec.Replicas; i++ {
			tCpuR.Add(cpuR)
			tCpuL.Add(cpuL)
			tMemR.Add(memR)
			tMemL.Add(memL)
			dataVol.Add(dv)
		}

	}

	aggrObj := Object{
		Spec: spec{
			DataVolumeCapacity: dataVol.String(),
			PodTemplate: &podTemplate{
				Resources: resources{
					Limits: res{
						Cpu:    tCpuL.String(),
						Memory: tMemL.String(),
					},
					Requests: res{
						Cpu:    tCpuR.String(),
						Memory: tMemR.String()},
				},
			},
		},
	}

	return aggrObj
}

func toOuputWithAggregates(objs []Object) (lines []OutputLine) {

	for _, o := range objs {

		if o.Spec.PodTemplate == nil {
			continue
		}

		line := OutputLine{
			Name:     o.Kind,
			Replicas: o.Spec.Replicas,
			CpuL:     o.Spec.PodTemplate.Resources.Limits.Cpu,
			MemoryL:  o.Spec.PodTemplate.Resources.Limits.Memory,
			CpuR:     o.Spec.PodTemplate.Resources.Requests.Cpu,
			MemoryR:  o.Spec.PodTemplate.Resources.Requests.Memory,
			DataVol:  o.Spec.DataVolumeCapacity,
		}

		lines = append(lines, line)
	}

	a := calculateAggregates(objs)
	aggrLine := OutputLine{
		CpuL:    a.Spec.PodTemplate.Resources.Limits.Cpu,
		MemoryL: a.Spec.PodTemplate.Resources.Limits.Memory,
		CpuR:    a.Spec.PodTemplate.Resources.Requests.Cpu,
		MemoryR: a.Spec.PodTemplate.Resources.Requests.Memory,
		DataVol: a.Spec.DataVolumeCapacity,
	}

	lines = append(lines, aggrLine)

	return
}

type Object struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Spec       spec   `yaml:"spec"`
}

type spec struct {
	DataVolumeCapacity string       `yaml:"dataVolumeCapacity"`
	PodTemplate        *podTemplate `yaml:"podTemplate"`
	Replicas           int          `yaml:"replicas"`
}

type podTemplate struct {
	Resources resources `yaml:"resources"`
}

type resources struct {
	Limits   res `yaml:"limits"`
	Requests res `yaml:"requests"`
}

type res struct {
	Cpu    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

type OutputLine struct {
	Name     string
	Replicas int
	CpuL     string
	MemoryL  string
	CpuR     string
	MemoryR  string
	DataVol  string
}

func (l OutputLine) String() string {
	return fmt.Sprintf("%s\t%d\t%s\t%s\t%s\t%s\t%s\t", l.Name, l.Replicas, l.CpuR, l.MemoryR, l.CpuL, l.MemoryL, l.DataVol)
}
