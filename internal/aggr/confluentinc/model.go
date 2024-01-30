package confluentinc

import "fmt"

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

func (o Object) Outputline() string {

	lf := "%s\t%d\t%s\t%s\t%s\t%s\t%s\t"
	name := o.Kind
	replicas := o.Spec.Replicas

	if o.Spec.PodTemplate == nil {
		return fmt.Sprintf(lf, name, replicas, "", "", "", "", o.Spec.DataVolumeCapacity)
	}

	return fmt.Sprintf(lf, name, replicas,
		o.Spec.PodTemplate.Resources.Requests.Cpu,
		o.Spec.PodTemplate.Resources.Requests.Memory,
		o.Spec.PodTemplate.Resources.Limits.Cpu,
		o.Spec.PodTemplate.Resources.Limits.Memory,
		o.Spec.DataVolumeCapacity)
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
