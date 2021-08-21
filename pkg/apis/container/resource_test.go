package container

import "testing"

func TestResourceRequirement_ToK8sResourceRequirement(t *testing.T) {
    r := ResourceRequirement{
        CPULimit:      100,
        MemoryLimit:   100,
        CPURequest:    95,
        MemoryRequest: 95,
    }
    k8sR := r.ToK8sResourceRequirement()
    if k8sR.Limits.Cpu().String() != "100Mi" {
        t.Error("cpu error")
    }
    if k8sR.Limits.Memory().String() != "100m" {
        t.Error("cpu error")
    }
    if k8sR.Requests.Cpu().String() != "95Mi" {
        t.Error("cpu error")
    }
    if k8sR.Requests.Memory().String() != "95m" {
        t.Error("cpu error")
    }
}
