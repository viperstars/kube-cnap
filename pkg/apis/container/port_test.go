package container

import "testing"

func TestPorts_ToK8sPort(t *testing.T) {
    ports := Ports{Ports: []*Port{{PortNumber: 8080}, {PortNumber: 9090}}}
    k8sPort := ports.ToK8sPort()
    if len(k8sPort) != 2 {
        t.Error("length error")
    }
    if k8sPort[0].ContainerPort != 8080 {
        t.Error("index error")
    }
}
