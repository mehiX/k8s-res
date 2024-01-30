package confluentinc

import (
	"k8s.io/apimachinery/pkg/api/resource"
)

func ComputeAggregates(objs Objects) Objects {

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

	aggrObj := newObject(
		"", "",
		tCpuL.String(),
		tMemL.String(),
		tCpuR.String(),
		tMemR.String(),
		dataVol.String(),
	)

	l := make([]Object, 0)
	l = append(l, objs...)
	l = append(l, aggrObj)

	return l
}
