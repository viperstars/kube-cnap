package pod

import (
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    v1 "k8s.io/api/core/v1"
    "regexp"
)

var separator = regexp.MustCompile(`\s*,\s*`)


type HostAlias struct {
    ID        int64             `json:"id" xorm:"pk autoincr 'id'"`
    Meta      basemeta.BaseMeta `json:"meta" xorm:"extends"`
    Hostnames map[string]string `json:"hostnames" xorm:""`
    Deleted   bool              `json:"deleted"`
}

func (h *HostAlias) ToK8sHostAliases() []v1.HostAlias {
    hosts := make([]v1.HostAlias, 0)
    for key, value := range h.Hostnames {
        hostnames := separator.Split(value, -1)
        host := v1.HostAlias{
            IP:        key,
            Hostnames: hostnames,
        }
        hosts = append(hosts, host)
    }
    return hosts
}
