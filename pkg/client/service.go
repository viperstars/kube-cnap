package client

import (
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    v1 "k8s.io/api/core/v1"
)

func (c *Clients) GetMiniServices(base basemeta.BaseMeta, srvList *v1.ServiceList) ([]*MiniService,
    error) {
    services := make([]*MiniService, 0)
    for _, s := range srvList.Items {
        miniService, _ := c.GetMiniServiceDetail(base, s, false)
        services = append(services, miniService)
    }
    return services, nil
}

func (c *Clients) GetMiniServiceDetail(base basemeta.BaseMeta, service v1.Service, detail bool) (*MiniService,
    error) {
    miniService := ToMiniService(service)
    if detail {
        lbs, _ := ToLabelSelector(service.Spec.Selector)
        pods, _ := c.GetPodsWithString(base, lbs.String())
        miniService.Pods = ToMiniPodList(pods.Items)
    }
    return miniService, nil
}
