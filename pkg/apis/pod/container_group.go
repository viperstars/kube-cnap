package pod

import (
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    "github.com/viperstars/kube-cnap/pkg/apis/consts"
    "github.com/viperstars/kube-cnap/pkg/apis/container"
    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ContainerGroup struct {
    ID                  int                    `json:"id" xorm:"pk autoincr 'id'"`
    Meta                basemeta.BaseMeta      `json:"meta" xorm:"extends"`
    MainContainerID     int                    `json:"mainContainer" xorm:"int 'main_container_id'" `
    InitContainersIDs   []int                  `json:"initContainers" xorm:"text 'init_containers_ids'"`
    SidecarContainerIDs []int                  `json:"sidecarContainers" xorm:"text 'sidecar_containers_ids'"`
    MainContainer       *container.Container   `json:"mainContainer" xorm:"-"`
    InitContainers      []*container.Container `json:"initContainers" xorm:"-"`
    SidecarContainers   []*container.Container `json:"sidecarContainers" xorm:"-"`
    NodeSelector        *NodeSelector          `json:"nodeSelector" xorm:"-"`
    HostAlias           *HostAlias           `json:"hostAlias" xorm:"-"`
    Tolerations         *Tolerations           `json:"tolerations" xorm:"-"`
    MainBase            *container.Base        `json:"mainBase" xorm:"-"`
    SidecarBases        []*container.Base      `json:"sidecarBases" xorm:"-"`
    InitBases           []*container.Base      `json:"initBases" xorm:"-"`
    Initialized         bool                   `json:"initialized" xorm:"bool default 1"`
    Deleted             bool                   `json:"deleted" xorm:"bool default 0"`
}

func (cg *ContainerGroup) ToNodeSelector() map[string]string {
    if cg.NodeSelector != nil {
        return cg.NodeSelector.ToK8sNodeSelector()
    } else {
        return nil
    }
}

func (cg *ContainerGroup) ToInitContainers() ([]v1.Container, []*v1.ConfigMap, []v1.Volume) {
    containers := make([]v1.Container, 0)
    configMaps := make([]*v1.ConfigMap, 0)
    volumes := make([]v1.Volume, 0)
    for _, _container := range cg.InitContainers {
        _container, configMap, volume := _container.ToK8sContainer(cg.Meta)
        containers = append(containers, *_container)
        configMaps = append(configMaps, configMap...)
        volumes = append(volumes, volume...)
    }
    return containers, configMaps, volumes
}

func (cg *ContainerGroup) ToContainers() ([]v1.Container, []*v1.ConfigMap, []v1.Volume) {
    containers := make([]v1.Container, 0)
    configMaps := make([]*v1.ConfigMap, 0)
    volumes := make([]v1.Volume, 0)
    _container, configMap, volume := cg.MainContainer.ToK8sContainer(cg.Meta)
    containers = append(containers, *_container)
    configMaps = append(configMaps, configMap...)
    volumes = append(volumes, volume...)
    for _, _container := range cg.SidecarContainers {
        _container, configMap, volume := _container.ToK8sContainer(cg.Meta)
        containers = append(containers, *_container)
        configMaps = append(configMaps, configMap...)
        volumes = append(volumes, volume...)
    }
    return containers, configMaps, volumes
}

func (cg *ContainerGroup) ToK8sPodTemplateSpecAndConfigMaps(deploymentName string, labels map[string]string,
    volumes []v1.Volume) (v1.PodTemplateSpec, []*v1.ConfigMap) {
    initContainers, configMapsForInitContainers, volumesForInitContainers := cg.ToInitContainers()
    containers, configMapsForContainers, volumesForContainers := cg.ToContainers()
    volumesForContainers = append(volumesForContainers, volumesForInitContainers...)
    configMapsForContainers = append(configMapsForContainers, configMapsForInitContainers...)
    volumesForContainers = append(volumesForContainers, volumes...)
    defaultSeconds := int64(60)
    podTemplateSpec := v1.PodTemplateSpec{
        ObjectMeta: metav1.ObjectMeta{
            Name:   deploymentName,
            Labels: labels,
        },
        Spec: v1.PodSpec{
            Volumes:                       volumesForContainers,
            InitContainers:                initContainers,
            Containers:                    containers,
            EphemeralContainers:           nil,
            RestartPolicy:                 "",
            TerminationGracePeriodSeconds: &defaultSeconds,
            ActiveDeadlineSeconds:         nil,
            DNSPolicy:                     "",
            NodeSelector:                  nil,
            ServiceAccountName:            "",
            DeprecatedServiceAccount:      "",
            AutomountServiceAccountToken:  nil,
            NodeName:                      "",
            HostNetwork:                   false,
            HostPID:                       false,
            HostIPC:                       false,
            ShareProcessNamespace:         nil,
            SecurityContext:               nil,
            ImagePullSecrets:              consts.DEFAULTSECRETNAME,
            Hostname:                      "",
            Subdomain:                     "",
            Affinity:                      nil,
            SchedulerName:                 "",
            Tolerations:                   nil,
            HostAliases:                   nil,
            PriorityClassName:             "",
            Priority:                      nil,
            DNSConfig:                     nil,
            ReadinessGates:                nil,
            RuntimeClassName:              nil,
            EnableServiceLinks:            nil,
            PreemptionPolicy:              nil,
            Overhead:                      nil,
            TopologySpreadConstraints:     nil,
        },
    }
    if cg.NodeSelector != nil && len(cg.NodeSelector.NodeSelector) > 0{
        podTemplateSpec.Spec.NodeSelector = cg.NodeSelector.ToK8sNodeSelector()
    }
    if cg.HostAlias != nil && len(cg.HostAlias.Hostnames) > 0 {
        podTemplateSpec.Spec.HostAliases = cg.HostAlias.ToK8sHostAliases()
    }
    if cg.Tolerations != nil {
        podTemplateSpec.Spec.Tolerations = cg.Tolerations.ToK8sTolerations()
    }
    return podTemplateSpec, configMapsForContainers
}
