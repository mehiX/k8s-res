package k8s

import (
	"k8s.io/apimachinery/pkg/api/resource"
)

func ComputeAggregates(objs Objects) Objects {

	tCpuR := new(resource.Quantity)
	tCpuL := new(resource.Quantity)
	tMemR := new(resource.Quantity)
	tMemL := new(resource.Quantity)

	for _, o := range objs {
		if o.IsEmpty() {
			continue
		}

		replicas := o.Spec.Replicas
		// missing replicas declaration, we assume 1 replica
		if replicas == 0 {
			replicas = 1
		}

		var containers []container
		containers = append(containers, o.Spec.Containers...)
		if o.Spec.Template != nil {
			containers = append(containers, o.Spec.Template.Spec.Containers...)
		}

		for _, container := range containers {

			var cpuR resource.Quantity
			if v := container.Resources.Requests.Cpu; v != "" {
				cpuR = resource.MustParse(v)
			}

			var memR resource.Quantity
			if v := container.Resources.Requests.Memory; v != "" {
				memR = resource.MustParse(v)
			}

			var cpuL resource.Quantity
			if v := container.Resources.Limits.Cpu; v != "" {
				cpuL = resource.MustParse(v)
			}

			var memL resource.Quantity
			if v := container.Resources.Limits.Memory; v != "" {
				memL = resource.MustParse(v)
			}

			for i := 0; i < replicas; i++ {
				tCpuR.Add(cpuR)
				tCpuL.Add(cpuL)
				tMemR.Add(memR)
				tMemL.Add(memL)
			}
		}
	}

	aggrObj := newObject(
		"", "",
		tCpuL.String(),
		tMemL.String(),
		tCpuR.String(),
		tMemR.String(),
	)

	l := make([]Object, 0)
	l = append(l, objs...)
	l = append(l, aggrObj)

	return l
}
