package response

import (
    "github.com/viperstars/kube-cnap/pkg/apis/ingress"
)

type HostIngress struct {
    Host    string `json:"host"`
    Ingress []*ingress.IngressRule
}

type IngressResponse struct {
    Ingress []*HostIngress
}

func ToIngressResponse(rules []*ingress.IngressRule) *IngressResponse {
    response := make([]*HostIngress, 0)
    sameHostPaths := make(map[string][]*ingress.IngressRule)
    for _, rule := range rules {
        if _, ok := sameHostPaths[rule.Host]; ok {
            sameHostPaths[rule.Host] = append(sameHostPaths[rule.Host], rule)
        } else {
            path := make([]*ingress.IngressRule, 0)
            sameHostPaths[rule.Host] = append(path, rule)
        }
    }
    for host, rules := range sameHostPaths {
        response = append(response, &HostIngress{
            Host:    host,
            Ingress: rules,
        })
    }
    return &IngressResponse{Ingress: response}
}
