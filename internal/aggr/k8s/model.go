package k8s

import (
	"fmt"
	"strings"
)

type Objects []Object

func newObject(apiVersion, kind string, cpuL, memL, cpuR, memR string) Object {

	return Object{
		Spec: spec{
			Template: &podTemplate{
				Spec: templateSpec{
					Containers: []container{{
						Resources: resources{
							Limits: res{
								Cpu:    cpuL,
								Memory: memL,
							},
							Requests: res{
								Cpu:    cpuR,
								Memory: memR},
						},
					}},
				},
			},
		},
	}
}

func (o Object) Outputline() string {

	lf := "%s\t%d\t%s\t%s\t%s\t%s\t"
	name := o.Kind
	replicas := o.Spec.Replicas

	if o.IsEmpty() {
		return fmt.Sprintf(lf, name, replicas, "", "", "", "")
	}

	var lines []string

	var containers []container
	containers = append(containers, o.Spec.Containers...)
	if o.Spec.Template != nil {
		containers = append(containers, o.Spec.Template.Spec.Containers...)
	}

	for _, container := range containers {
		line := fmt.Sprintf(lf, name, replicas,
			container.Resources.Requests.Cpu,
			container.Resources.Requests.Memory,
			container.Resources.Limits.Cpu,
			container.Resources.Limits.Memory)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (o Object) IsEmpty() bool {
	return o.Spec.Template == nil && len(o.Spec.Containers) == 0
}

type Object struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Spec       spec   `yaml:"spec"`
}

type spec struct {
	Template   *podTemplate `yaml:"template"`
	Replicas   int          `yaml:"replicas"`
	Containers []container  `yaml:"containers"`
}

type podTemplate struct {
	Spec templateSpec `yaml:"spec"`
}

type templateSpec struct {
	Containers []container `yaml:"containers"`
}

type container struct {
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
