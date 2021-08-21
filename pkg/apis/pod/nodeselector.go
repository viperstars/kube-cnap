package pod

import "github.com/viperstars/kube-cnap/pkg/apis/basemeta"

type NodeSelector struct {
    ID           int64             `json:"id" xorm:"pk autoincr 'id'"`
    Meta         basemeta.BaseMeta `json:"meta" xorm:"extends"`
    NodeSelector map[string]string `json:"nodeSelector"`
    Deleted      bool     `json:"deleted"`
}

func (n *NodeSelector) ToK8sNodeSelector() map[string]string {
    return n.NodeSelector
}
