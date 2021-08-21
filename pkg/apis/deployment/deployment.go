package deployment

import (
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    "github.com/viperstars/kube-cnap/pkg/apis/consts"
    appsv1 "k8s.io/api/apps/v1"
    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/util/intstr"
)

type Deployment struct {
    Number         int32  `json:"number"`
    Name           string `json:"name"`
    Image          string `json:"image"`
    MaxUnavailable int32 `json:"maxUnavailable"`
    MaxSurge       int32 `json:"maxSurge"`
}

func (d *Deployment) GetLabels(base basemeta.BaseMeta) map[string]string {
    labels := base.GetLabels()
    labels["deployment"] = d.Name
    return labels
}

func (d *Deployment) ToK8sStrategy() appsv1.DeploymentStrategy {
    strategy := appsv1.DeploymentStrategy{
        Type: consts.DEFAULTROLLINGUPDATESTRATEGYTYPE,
        RollingUpdate: &appsv1.RollingUpdateDeployment{
            MaxUnavailable: &intstr.IntOrString{
                IntVal: d.MaxUnavailable,
            },
            MaxSurge: &intstr.IntOrString{
                IntVal: d.MaxSurge,
            },
        },
    }
    return strategy
}

func (d *Deployment) ToK8sDeployment(base basemeta.BaseMeta, podTemplateSpec v1.PodTemplateSpec) *appsv1.Deployment {
    labels := d.GetLabels(base)
    k8sDeployment := &appsv1.Deployment{
        TypeMeta: metav1.TypeMeta{},
        ObjectMeta: metav1.ObjectMeta{
            Name:   d.Name,
            Labels: labels,
        },
        Spec: appsv1.DeploymentSpec{
            Replicas: &d.Number,
            Template: podTemplateSpec,
            Strategy: d.ToK8sStrategy(),
            Selector: &metav1.LabelSelector{
                MatchLabels: labels,
            },
            MinReadySeconds:         0,
            RevisionHistoryLimit:    nil,
            Paused:                  false,
            ProgressDeadlineSeconds: nil,
        },
        Status: appsv1.DeploymentStatus{},
    }
    return k8sDeployment
}
