package pod

import (
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    v1 "k8s.io/api/core/v1"
)

type Toleration struct {
    ID       int               `json:"id" xorm:"pk autoincr 'id'"`
    Meta     basemeta.BaseMeta `json:"meta" xorm:"extends"`
    Key      string            `json:"key"`
    Value    string            `json:"value"`
    Operator string            `json:"operator"`
    Effect   string            `json:"effect"`
    Seconds  int               `json:"seconds"`
}

type Tolerations struct {
    Tolerations []*Toleration
}

func (t *Tolerations) ToK8sTolerations() []v1.Toleration {
    ts := make([]v1.Toleration, 0)
    for _, toleration := range t.Tolerations {
        tt := v1.Toleration{
            Key:      toleration.Key,
            Value:    toleration.Value,
            Effect:   v1.TaintEffect(toleration.Effect),
            Operator: v1.TolerationOperator(toleration.Operator),
        }
        ts = append(ts, tt)
    }
    return ts
}
