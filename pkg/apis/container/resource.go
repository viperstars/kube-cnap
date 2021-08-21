package container

import (
    "fmt"
    v1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/resource"
)

type ResourceRequirement struct {
    ID            int    `json:"id" xorm:"pk autoincr 'id'"`
    ContainerID   int    `json:"containerID" xorm:"int 'container_id'"`
    CPULimit      int    `json:"cpuLimit" xorm:"int 'cpu_limit'"`
    MemoryLimit   int    `json:"memoryLimit" xorm:"int"`
    CPURequest    int    `json:"cpuRequest" xorm:"int 'cpu_request'"`
    MemoryRequest int    `json:"memoryRequest" xorm:"int 'memory_request'"`
    Deleted       bool   `json:"deleted" xorm:"bool default 0"`
}

func (resourceRequirement *ResourceRequirement) ToK8sResourceRequirement() v1.ResourceRequirements {
    limit := v1.ResourceList{}
    request := v1.ResourceList{}
    if resourceRequirement.CPULimit > 0 {
        value := fmt.Sprintf("%dm", resourceRequirement.CPULimit)
        limit[v1.ResourceCPU] = resource.MustParse(value)
    }
    if resourceRequirement.MemoryLimit > 0 {
        value := fmt.Sprintf("%dMi", resourceRequirement.MemoryLimit)
        limit[v1.ResourceMemory] = resource.MustParse(value)
    }
    if resourceRequirement.CPURequest > 0 {
        value := fmt.Sprintf("%dm", resourceRequirement.CPURequest )
        request[v1.ResourceCPU] = resource.MustParse(value)
    }
    if resourceRequirement.MemoryRequest > 0 {
        value := fmt.Sprintf("%dMi", resourceRequirement.MemoryRequest)
        request[v1.ResourceMemory] = resource.MustParse(value)
    }
    k8sResourceRequirement := v1.ResourceRequirements{
        Limits:   limit,
        Requests: request,
    }
    return k8sResourceRequirement
}

// todo resource 单位的问题
// todo 属性为空时的处理逻辑