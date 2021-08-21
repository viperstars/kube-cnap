package ingress

import (
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    "github.com/viperstars/kube-cnap/pkg/naming"
    "k8s.io/api/extensions/v1beta1"
    v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/util/intstr"
)

type IngressRules struct {
    IngressRules []*IngressRule
}

type IngressRule struct {
    ID          int               `json:"id" xorm:"pk autoincr 'id'"`
    Meta        basemeta.BaseMeta `json:"meta" xorm:"extends"`
    Host        string            `json:"host" xorm:"varchar(512)"`
    Path        string            `json:"path" xorm:"varchar(512)"`
    ServiceName string            `json:"serviceName" xorm:"varchar(256)"`
    ServicePort int32             `json:"servicePort" xorm:"varchar(128)"`
    Deleted     bool              `json:"deleted"`
}

func (i *IngressRules) GetSameHostPaths() map[string][]v1beta1.HTTPIngressPath {
    sameHostPaths := make(map[string][]v1beta1.HTTPIngressPath)
    for _, rule := range i.IngressRules {
        if _, ok := sameHostPaths[rule.Host]; ok {
            sameHostPaths[rule.Host] = append(sameHostPaths[rule.Host], rule.ToK8sIngressRule()...)
        } else {
            path := make([]v1beta1.HTTPIngressPath, 0)
            sameHostPaths[rule.Host] = append(path, rule.ToK8sIngressRule()...)
        }
    }
    return sameHostPaths
}

func (i *IngressRules) ToK8sIngressRules() []v1beta1.IngressRule {
    paths := i.GetSameHostPaths()
    rules := make([]v1beta1.IngressRule, 0)
    for host, path := range paths {
        ps := make([]v1beta1.HTTPIngressPath, 0)
        for _, p := range path {
            ps = append(ps, p)
        }
        r := v1beta1.IngressRule{
            Host: host,
            IngressRuleValue: v1beta1.IngressRuleValue{
                HTTP: &v1beta1.HTTPIngressRuleValue{
                    Paths: ps,
                },
            },
        }
        rules = append(rules, r)
    }
    return rules
}

func (i *IngressRule) ToK8sBackend() v1beta1.IngressBackend {
    backend := v1beta1.IngressBackend{
        ServiceName: i.ServiceName,
        ServicePort: intstr.IntOrString{
            IntVal: i.ServicePort,
        },
    }
    return backend
}

func (i *IngressRule) ToK8sIngressRule() []v1beta1.HTTPIngressPath {
    paths := []v1beta1.HTTPIngressPath{
        {i.Path, i.ToK8sBackend()},
    }
    return paths
}

func (i *IngressRules) ToK8sIngress(base basemeta.BaseMeta) *v1beta1.Ingress {
    rules := i.ToK8sIngressRules()
    ingress := &v1beta1.Ingress{
        TypeMeta: v1.TypeMeta{},
        ObjectMeta: v1.ObjectMeta{
            Name: naming.Namer.IngressName(&base),
        },
        Spec: v1beta1.IngressSpec{
            TLS:   nil,
            Rules: rules,
        },
    }
    return ingress
}
