package container

import v1 "k8s.io/api/core/v1"

type Ports struct {
    Ports []*Port `json:"ports"`
}

type Port struct {
    ID          int    `json:"id" xorm:"pk autoincr 'id'"`
    ContainerID int    `json:"containerID" xorm:"int 'container_id'"`
    PortName    string `json:"portName" xorm:"varchar(256)"`
    PortNumber  int32  `json:"portNumber" xorm:"int"`
    Protocol    string `json:"protocol"`
    Deleted     bool   `json:"deleted" xorm:"bool default 0"`
}

func (p *Ports) ToK8sPort() []v1.ContainerPort {
    containerPort := make([]v1.ContainerPort, 0)
    for _, port := range p.Ports {
        var protocol v1.Protocol
        switch port.Protocol {
        case string(v1.ProtocolSCTP):
            protocol = v1.ProtocolSCTP
        case string(v1.ProtocolUDP):
            protocol = v1.ProtocolUDP
        default:
            protocol = v1.ProtocolTCP
        }
        cp := v1.ContainerPort{
            Name: port.PortName,
            ContainerPort: port.PortNumber,
            Protocol: protocol,
        }
        containerPort = append(containerPort, cp)
    }
    return containerPort
}
