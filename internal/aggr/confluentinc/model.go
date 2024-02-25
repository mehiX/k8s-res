package confluentinc

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"
)

type Objects []Object

func newObject(apiVersion, kind string, cpuL, memL, cpuR, memR, dataVol string) Object {

	return Object{
		Spec: spec{
			DataVolumeCapacity: dataVol,
			PodTemplate: &podTemplate{
				Resources: resources{
					Limits: res{
						Cpu:    cpuL,
						Memory: memL,
					},
					Requests: res{
						Cpu:    cpuR,
						Memory: memR},
				},
			},
		},
	}
}

func (o Object) diskVol() string {
	diskVol := new(resource.Quantity)
	if v := o.Spec.DataVolumeCapacity; v != "" {
		diskVol.Add(resource.MustParse(v))
	}
	if v := o.Spec.LogVolumeCapacity; v != "" {
		diskVol.Add(resource.MustParse(v))
	}

	return diskVol.String()
}

func (o Object) String() string {

	lf := "%s\t%s\t%s\t%s\t%s\t%s\t%s\t"
	name := o.Kind
	replicas := ""
	if o.Spec.Replicas > 0 {
		replicas = fmt.Sprintf("%d", o.Spec.Replicas)
	}

	if o.Spec.PodTemplate == nil {
		return fmt.Sprintf(lf, name, replicas, "", "", "", "", o.diskVol())
	}

	return fmt.Sprintf(lf, name, replicas,
		o.Spec.PodTemplate.Resources.Requests.Cpu,
		o.Spec.PodTemplate.Resources.Requests.Memory,
		o.Spec.PodTemplate.Resources.Limits.Cpu,
		o.Spec.PodTemplate.Resources.Limits.Memory,
		o.diskVol())
}

func (o Object) IsEmpty() bool {
	return o.Spec.PodTemplate == nil
}

type Object struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Spec       spec   `yaml:"spec"`
}

type spec struct {
	DataVolumeCapacity string       `yaml:"dataVolumeCapacity"`
	LogVolumeCapacity  string       `yaml:"logVolumeCapacity"`
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
