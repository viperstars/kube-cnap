package client

import (
    "fmt"
    "github.com/viperstars/kube-cnap/pkg/apis/consts"
    "github.com/viperstars/kube-cnap/pkg/apis/service"
    appsv1 "k8s.io/api/apps/v1"
    v1 "k8s.io/api/core/v1"
)

type PodInfo struct {
    Current   int32 `json:"current"`
    Desired   int32 `json:"desired,omitempty"`
    Running   int32 `json:"running"`
    Pending   int32 `json:"pending"`
    Failed    int32 `json:"failed"`
    Succeeded int32 `json:"succeeded"`
}

type MiniStatus struct {
    Name         string `json:"name"`
    RestartCount int32  `json:"restartCount"`
    Ready        bool   `json:"ready"`
    States       *v1.ContainerStateTerminated `json:"states"`
}

type MiniPod struct {
    Name      string        `json:"name"`
    NodeName  string        `json:"nodeName"`
    PodIP     string        `json:"podIP"`
    IP        string        `json:"ip"`
    StartTime string        `json:"startTime"`
    QoSClass  string        `json:"qosClass"`
    Status    []*MiniStatus `json:"status"`
    Type      string        `json:"type"`
}

func GetPodInfo(current int32, desired int32, pods []v1.Pod) PodInfo {
    result := PodInfo{
        Current: current,
        Desired: desired,
    }
    for _, p := range pods {
        switch p.Status.Phase {
        case v1.PodRunning:
            result.Running++
        case v1.PodPending:
            result.Pending++
        case v1.PodFailed:
            result.Failed++
        case v1.PodSucceeded:
            result.Succeeded++
        }
    }

    return result
}

type ResourceStatus struct {
    Running   int `json:"running"`
    Pending   int `json:"pending"`
    Failed    int `json:"failed"`
    Succeeded int `json:"succeeded"`
}

type SimpleDeployment struct {
    Deployment  *appsv1.Deployment  `json:"deployment"`
    ReplicaSets []*SimpleReplicaSet `json:"replicaSets"`
    PodInfo     PodInfo             `json:"podInfo"`
}

type SimpleReplicaSet struct {
    ReplicaSet *appsv1.ReplicaSet `json:"replicaSet"`
    PodList    []v1.Pod           `json:"podList"`
    PodInfo    PodInfo            `json:"podInfo"`
}

type SimpleDeploymentList struct {
    Deployments []*SimpleDeployment `json:"deployments"`
}

func ToMiniStatus(status []v1.ContainerStatus) []*MiniStatus {
    miniStatus := make([]*MiniStatus, 0)
    for _, st := range status {
        mini := new(MiniStatus)
        mini.Name = st.Name
        mini.RestartCount = st.RestartCount
        mini.Ready = st.Ready
        mini.States = st.LastTerminationState.Terminated
        miniStatus = append(miniStatus, mini)
    }
    return miniStatus
}

func ToMiniPod(pod v1.Pod) *MiniPod {
    miniPod := new(MiniPod)
    miniPod.Name = pod.Name
    miniPod.NodeName = pod.Spec.NodeName
    if pod.Status.StartTime != nil {
        miniPod.StartTime = pod.Status.StartTime.String()
    }
    miniPod.QoSClass = string(pod.Status.QOSClass)
    miniPod.PodIP = pod.Status.PodIP
    miniPod.IP = pod.Status.PodIP
    miniPod.Status = ToMiniStatus(pod.Status.ContainerStatuses)
    miniPod.Type = "容器组"
    return miniPod
}

type MiniDeployment struct {
    Name              string            `json:"name"`
    Replicas          int32             `json:"replicas"`
    UpdatedReplicas   int32             `json:"updatedReplicas"`
    ReadyReplicas     int32             `json:"readyReplicas"`
    AvailableReplicas int32             `json:"availableReplicas"`
    Strategy          string            `json:"strategy"`
    MaxUnavailable    int32             `json:"maxUnavailable"`
    MaxSurge          int32             `json:"maxSurge"`
    NewestRs          *MiniReplicaSet   `json:"newestRs"`
    OtherRs           []*MiniReplicaSet `json:"otherRs"`
    IP                string            `json:"ip"`
    Children          interface{}       `json:"children"`
    Type              string            `json:"type"`
}

func ToMiniDeployment(deployment appsv1.Deployment) *MiniDeployment {
    miniDeployment := new(MiniDeployment)
    miniDeployment.Type = "集群"
    miniDeployment.Name = deployment.Name
    miniDeployment.Replicas = *deployment.Spec.Replicas
    miniDeployment.Strategy = string(deployment.Spec.Strategy.Type)
    miniDeployment.MaxSurge = deployment.Spec.Strategy.RollingUpdate.MaxSurge.IntVal
    miniDeployment.MaxUnavailable = deployment.Spec.Strategy.RollingUpdate.MaxUnavailable.IntVal
    miniDeployment.AvailableReplicas = deployment.Status.AvailableReplicas
    miniDeployment.ReadyReplicas = deployment.Status.ReadyReplicas
    miniDeployment.UpdatedReplicas = deployment.Status.UpdatedReplicas
    miniDeployment.IP = ""
    return miniDeployment
}

func ToMiniPodList(list []v1.Pod) []*MiniPod {
    podList := make([]*MiniPod, 0)
    for _, p := range list {
        miniPod := ToMiniPod(p)
        podList = append(podList, miniPod)
    }
    return podList
}

func GetNewReplicaSetTemplate(deployment appsv1.Deployment) v1.PodTemplateSpec {
    return v1.PodTemplateSpec{
        ObjectMeta: deployment.Spec.Template.ObjectMeta,
        Spec:       deployment.Spec.Template.Spec,
    }
}

type MiniReplicaSet struct {
    Name    string     `json:"name"`
    Image   string     `json:"image"`
    PodList []*MiniPod `json:"podList"`
    PodInfo *PodInfo   `json:"podInfo"`
}

type MiniDeploymentList struct {
    DeploymentList []*MiniDeployment `json:"deploymentList"`
}

func FilterReplicaSet(deployment appsv1.Deployment, rsList []appsv1.ReplicaSet) (appsv1.ReplicaSet,
    []appsv1.ReplicaSet) {
    var newest appsv1.ReplicaSet
    others := make([]appsv1.ReplicaSet, 0)
    newRSTemplate := GetNewReplicaSetTemplate(deployment)
    for i := range rsList {
        if EqualIgnoreHash(rsList[i].Spec.Template, newRSTemplate) {
            newest = rsList[i]
        } else {
            others = append(others, rsList[i])
        }
    }
    return newest, others
}

type MiniService struct {
    Name              string                  `json:"name"`
    CreationTimestamp string                  `json:"creationTimestamp"`
    ClusterIP         string                  `json:"clusterIP"`
    SessionAffinity   string                  `json:"sessionAffinity"`
    Ports             []*service.ResponsePort `json:"ports"`
    Pods              []*MiniPod              `json:"pods"`
    DeploymentName    string                  `json:"deploymentName"`
    ExternalName      string                  `json:"externalName"`
    Comment           string                  `json:"comment"`
    IP                string                  `json:"ip"`
    Children          interface{}             `json:"children"`
    Type              string                  `json:"type"`
    ServiceType       string                  `json:"serviceType"`
    StickySession     bool                    `json:"stickySession"`
}

func ToMiniService(service v1.Service) *MiniService {
    miniService := new(MiniService)
    miniService.Type = "服务"
    miniService.Name = service.Name
    miniService.CreationTimestamp = service.CreationTimestamp.String()
    miniService.ClusterIP = service.Spec.ClusterIP
    miniService.SessionAffinity = string(service.Spec.SessionAffinity)
    miniService.Ports = ToMiniServicePorts(service.Spec.Ports)
    miniService.ServiceType = string(service.Spec.Type)
    if _, ok := service.Spec.Selector["deployment"]; ok {
        miniService.DeploymentName = service.Spec.Selector["deployment"]
    } else {
        miniService.DeploymentName = "all"
    }
    if _, ok := service.Annotations[consts.DEFAULTSTICKYSESSIONKEY]; ok {
        miniService.StickySession = true
    }
    if _, ok := service.Annotations["externalName"]; ok {
        miniService.ExternalName = service.Annotations["externalName"]
    } else {
        miniService.ExternalName = ""
    }
    if _, ok := service.Annotations["comment"]; ok {
        miniService.Comment = service.Annotations["comment"]
    } else {
        miniService.Comment = ""
    }
    return miniService
}

func ToMiniServicePorts(servicePorts []v1.ServicePort) []*service.ResponsePort {
    ports := make([]*service.ResponsePort, 0)
    for _, p := range servicePorts {
        port := &service.ResponsePort{
            Name:       p.Name,
            Protocol:   string(p.Protocol),
            Port:       p.Port,
            TargetPort: p.TargetPort.IntVal,
            NodePort:   p.NodePort,
        }
        ports = append(ports, port)
    }
    return ports
}

type MiniEvent struct {
    Message         string `json:"message"`
    SourceComponent string `json:"sourceComponent"`
    SourceHost      string `json:"sourceHost"`
    Object          string `json:"object"`
    Count           int32  `json:"count"`
    FirstSeen       string `json:"firstSeen"`
    LastSeen        string `json:"lastSeen"`
    Reason          string `json:"reason"`
    Type            string `json:"type"`
}

func ToMiniEvent(event v1.Event) *MiniEvent {
    result := &MiniEvent{
        Message:         event.Message,
        SourceComponent: event.Source.Component,
        SourceHost:      event.Source.Host,
        Object:          event.InvolvedObject.FieldPath,
        Count:           event.Count,
        FirstSeen:       event.FirstTimestamp.String(),
        LastSeen:        event.LastTimestamp.String(),
        Reason:          event.Reason,
        Type:            event.Type,
    }
    fmt.Println(event)
    return result
}

func ToMiniEventList(list v1.EventList) []*MiniEvent {
    events := make([]*MiniEvent, 0)
    for _, e := range list.Items {
        event := ToMiniEvent(e)
        events = append(events, event)
    }
    return events
}

type MiniNode struct {
    Name                    string            `json:"name"`
    Labels                  map[string]string `json:"labels"`
    AllocatableCPU          string            `json:"allocatableCPU"`
    AllocatableMemory       int64            `json:"allocatableMemory"`
    AllocatablePods         string            `json:"allocatablePods"`
    TotalCPU                string            `json:"totalCPU"`
    TotalMemory             int64            `json:"totalMemory"`
    TotalPods               string            `json:"totalPods"`
    IP                      string            `json:"ip"`
    KernelVersion           string            `json:"kernelVersion"`
    OSImage                 string            `json:"osImage"`
    ContainerRuntimeVersion string            `json:"containerRuntimeVersion"`
    KubeletVersion          string            `json:"kubeletVersion"`
    KubeProxyVersion        string            `json:"kubeProxyVersion"`
    OperatingSystem         string            `json:"operatingSystem"`
    Architecture            string            `json:"architecture"`
}

func ToMiniNode(node v1.Node) *MiniNode {
    miniNode := &MiniNode{
        Name:                    node.Name,
        Labels:                  node.Labels,
        AllocatableCPU:          node.Status.Allocatable.Cpu().String(),
        AllocatableMemory:       node.Status.Allocatable.Memory().Value() / (1024 * 1024),
        AllocatablePods:         node.Status.Allocatable.Pods().String(),
        TotalCPU:                node.Status.Capacity.Cpu().String(),
        TotalMemory:             node.Status.Capacity.Memory().Value() / (1024 * 1024),
        TotalPods:               node.Status.Capacity.Pods().String(),
        IP:                      node.Status.Addresses[0].Address,
        KernelVersion:           node.Status.NodeInfo.KernelVersion,
        OSImage:                 node.Status.NodeInfo.OSImage,
        ContainerRuntimeVersion: node.Status.NodeInfo.ContainerRuntimeVersion,
        KubeletVersion:          node.Status.NodeInfo.KubeletVersion,
        KubeProxyVersion:        node.Status.NodeInfo.KubeProxyVersion,
        OperatingSystem:         node.Status.NodeInfo.OperatingSystem,
        Architecture:            node.Status.NodeInfo.Architecture,
    }
    return miniNode
}

func ToMiniNodeList(list v1.NodeList) []*MiniNode {
    nodes := make([]*MiniNode, 0)
    for _, n := range list.Items {
        node := ToMiniNode(n)
        nodes = append(nodes, node)
    }
    return nodes
}

func ToMiniServiceList(list v1.ServiceList) []*MiniService {
    services := make([]*MiniService, 0)
    for _, s := range list.Items {
        srv := ToMiniService(s)
        services = append(services, srv)
    }
    return services
}
