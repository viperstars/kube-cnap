package resourcequota

import (
    "fmt"
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    corev1 "k8s.io/api/core/v1"
    v1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/resource"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ResourceQuota struct {
    Base          basemeta.BaseMeta `json:"base"`
    RequestCPU    int64 `json:"requestCpu"`
    RequestMemory int64 `json:"requestMemory"`
    LimitCPU      int64 `json:"limitCpu"`
    LimitMemory   int64 `json:"limitMemory"`
    Region        string `json:"region"`
    Env           string `json:"env"`
}

func (r *ResourceQuota) ToK8sResourceQuota() *v1.ResourceQuota {
    hard := make(map[corev1.ResourceName]resource.Quantity)
    requestCpu := fmt.Sprintf("%d", r.RequestCPU)
    hard["requests.cpu"] = resource.MustParse(requestCpu)
    limitCpu := fmt.Sprintf("%d", r.LimitCPU)
    hard["limits.cpu"] = resource.MustParse(limitCpu)
    requestMemory := fmt.Sprintf("%dGi", r.RequestMemory)
    hard["requests.memory"] = resource.MustParse(requestMemory)
    limitMemory := fmt.Sprintf("%dGi", r.LimitMemory)
    hard["limits.memory"] = resource.MustParse(limitMemory)
    rq := &v1.ResourceQuota{
        TypeMeta:   metav1.TypeMeta{},
        ObjectMeta: metav1.ObjectMeta{
            Name: r.Base.App,
        },
        Spec: v1.ResourceQuotaSpec{
            Hard: hard,
        },
        Status:     v1.ResourceQuotaStatus{},
    }
    return rq
}
