package service

import (
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    "github.com/viperstars/kube-cnap/pkg/apis/consts"
    "github.com/viperstars/kube-cnap/pkg/naming"
    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/util/intstr"
)

type Port struct {
    Name       string `json:"name"`
    Protocol   string `json:"protocol"`
    Port       int32  `json:"port"`
    TargetPort int32  `json:"targetPort"`
}

type ResponsePort struct {
    Name       string `json:"name"`
    Protocol   string `json:"protocol"`
    Port       int32  `json:"port"`
    TargetPort int32  `json:"targetPort"`
    NodePort   int32  `json:"nodePort"`
}

func (p *Port) ToK8sServicePort() v1.ServicePort {
    var protocol v1.Protocol
    switch p.Protocol {
    case string(v1.ProtocolSCTP):
        protocol = v1.ProtocolSCTP
    case string(v1.ProtocolUDP):
        protocol = v1.ProtocolUDP
    default:
        protocol = v1.ProtocolTCP
    }
    return v1.ServicePort{
        Protocol: protocol,
        Port:     p.Port,
        TargetPort: intstr.IntOrString{
            IntVal: p.TargetPort,
        },
    }
}

type CreateServiceRequest struct {
    Name           string  `json:"name"`
    Ports          []*Port `json:"ports"`
    DeploymentName string  `json:"deploymentName"`
    ServiceType    string  `json:"serviceType"`
    ExternalName  string `json:"externalName"`
    StickySession bool   `json:"stickySession"`
}

func (c *CreateServiceRequest) GetFullName(base basemeta.BaseMeta) string {
    return naming.Namer.ServiceName(&base, c.Name)
}

func (c *CreateServiceRequest) ToLabelSelector(base basemeta.BaseMeta) map[string]string { // service only support
    baseLabels := base.GetLabels()
    if c.DeploymentName != "all" {
        baseLabels["deployment"] = c.DeploymentName
    }
    return baseLabels
}

func (c *CreateServiceRequest) ToK8sService(base basemeta.BaseMeta) *v1.Service { // service only support
    var t v1.ServiceType
    ports := make([]v1.ServicePort, 0)
    for _, p := range c.Ports {
        ports = append(ports, p.ToK8sServicePort())
    }
    switch c.ServiceType {
    case "NodePort":
        t = v1.ServiceTypeNodePort
    case "LoadBalancer":
        t = v1.ServiceTypeLoadBalancer
    default:
        t = v1.ServiceTypeClusterIP
    }
    annotation := make(map[string]string, 0)
    if c.ExternalName != "" {
        annotation["externalName"] = c.ExternalName
    }
    if c.StickySession == true {
        annotation[consts.DEFAULTSTICKYSESSIONKEY] = "true"
    }
    service := &v1.Service{
        TypeMeta: metav1.TypeMeta{},
        ObjectMeta: metav1.ObjectMeta{
            Name:   c.Name,
            Labels: base.GetLabels(),
            Annotations: annotation,
        },
        Spec: v1.ServiceSpec{
            Ports:           ports,
            Type:            t,
            Selector:        c.ToLabelSelector(base),
            SessionAffinity: consts.DEFAULTSERVICEAFFINITY,
        },
        Status: v1.ServiceStatus{},
    }
    return service
}
