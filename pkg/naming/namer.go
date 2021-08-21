package naming

import (
    "crypto/md5"
    "fmt"
)

var Namer *Naming

func init() {
    Namer = new(Naming)
}

type Prefix interface {
    GetPrefix() string
}

type Naming struct{}

func (n Naming) ConfigMapName(prefix Prefix, dirname string) string {
    return n.getPathHash(dirname)
}

func (n Naming) getPathHash(path string) string {
    return fmt.Sprintf("%x", md5.Sum([]byte(path)))
}

func (n Naming) VolumeName(prefix string, volumeName string) string {
    return prefix + "-" + volumeName
}

func (n Naming) DeploymentName(prefix Prefix, deploymentName string) string {
    return prefix.GetPrefix() + "-" + deploymentName
}

func (n Naming) ServiceName(prefix Prefix, serviceName string) string {
    return prefix.GetPrefix() + "-" + serviceName
}

func (n Naming) IngressName(prefix Prefix) string {
    return prefix.GetPrefix() + "-" + "ingress"
}
