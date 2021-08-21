package requests

import (
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    "github.com/viperstars/kube-cnap/pkg/apis/consts"
    appsv1 "k8s.io/api/apps/v1"
    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/util/intstr"
)

type CreateDeploymentRequest struct {
    Number         int32  `json:"number"`
    Name           string `json:"name"`
    Image          string `json:"image"`
    MaxUnavailable int32 `json:"maxUnavailable"`
    MaxSurge       int32 `json:"maxSurge"`
}

func (r *CreateDeploymentRequest) GetLabels(base basemeta.BaseMeta) map[string]string {
    lbs := make(map[string]string)
    lbs["app"] = base.App
    lbs["region"] = base.Region
    lbs["env"] = base.Env
    lbs["deployment"] = r.Name
    return lbs
}

func (r *CreateDeploymentRequest) ToK8sStrategy() appsv1.DeploymentStrategy {
    strategy := appsv1.DeploymentStrategy{
        Type: consts.DEFAULTROLLINGUPDATESTRATEGYTYPE,
        RollingUpdate: &appsv1.RollingUpdateDeployment{
            MaxUnavailable: &intstr.IntOrString{
                IntVal: r.MaxUnavailable,
            },
            MaxSurge: &intstr.IntOrString{
                IntVal: r.MaxSurge,
            },
        },
    }
    return strategy
}

func (r *CreateDeploymentRequest) ToK8sDeployment(lbs map[string]string, podTemplateSpec v1.PodTemplateSpec) *appsv1.Deployment {
    k8sDeployment := &appsv1.Deployment{
        TypeMeta: metav1.TypeMeta{},
        ObjectMeta: metav1.ObjectMeta{
            Name:   r.Name,
            Labels: lbs,
        },
        Spec: appsv1.DeploymentSpec{
            Replicas: &r.Number,
            Template: podTemplateSpec,
            Strategy: r.ToK8sStrategy(),
            Selector: &metav1.LabelSelector{
                MatchLabels: lbs,
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

type UpdateDeploymentRequest struct {
    Image string `json:"image"`
    VersionOnly bool `json:"versionOnly"`
}

type ScaleDeploymentRequest struct {
    Replicas int `json:"replicas"`
}
