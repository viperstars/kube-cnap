package client

import (
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
)

func (c *Clients) GetMiniDeploymentList(base basemeta.BaseMeta) ([]*MiniDeployment, error) {
    deployments, err := c.GetDeployments(base, base.GetLabels())
    if err != nil {
        return nil, err
    }
    return c.GetMiniDeployments(base, deployments)
}

func (c *Clients) GetMiniDeployment(base basemeta.BaseMeta, deploymentName string) (*MiniDeployment, error) {
    deployment, err := c.GetDeployment(base, deploymentName)
    if err != nil {
        return nil, err
    }
    return c.GetMiniDeploymentDetail(base, *deployment, true)
}

func (c *Clients) GetMiniServiceList(base basemeta.BaseMeta) ([]*MiniService, error) {
    services, err := c.GetServices(base, base.GetLabels())
    if err != nil {
        return nil, err
    }
    return c.GetMiniServices(base, services)
}

func (c *Clients) GetMiniService(base basemeta.BaseMeta, serviceName string) (*MiniService, error) {
    service, err := c.GetService(base, serviceName)
    if err != nil {
        return nil, err
    }
    return c.GetMiniServiceDetail(base, *service, true)
}

func (c *Clients) GetMiniEventList(base basemeta.BaseMeta, objectKind string, objectName string) ([]*MiniEvent, error) {
    var events []*MiniEvent
    var err error
    if objectKind != "node" {
        events, err = c.GetMiniEventListForObject(base, objectName)
    } else {
        events, err = c.GetMiniEventListForNode(base, objectName, true)
    }
    if err != nil {
        return nil, err
    }
    return events, nil
}

func (c *Clients) GetMiniEventListForObject(base basemeta.BaseMeta, objectName string) ([]*MiniEvent,
    error) {
    events, err := c.GetEventsForObject(base, objectName)
    if err != nil {
        return nil, err
    }
    return ToMiniEventList(*events), nil
}

func (c *Clients) GetMiniEventListForNode(base basemeta.BaseMeta, nodeName string, allNamespaces bool) ([]*MiniEvent,
    error) {
    events, err := c.GetEventsForNode(base, nodeName, allNamespaces)
    if err != nil {
        return nil, err
    }
    return ToMiniEventList(*events), nil
}

func (c *Clients) GetMiniNodesList(base basemeta.BaseMeta) ([]*MiniNode, error) {
    nodes, err := c.GetNodes(base)
    if err != nil {
        return nil, err
    }
    return ToMiniNodeList(*nodes), nil
}

func (c *Clients) GetMiniNode(base basemeta.BaseMeta, nodeName string) (*MiniNode, error) {
    node, err := c.GetNode(base, nodeName)
    if err != nil {
        return nil, err
    }
    return ToMiniNode(*node), nil
}

func (c *Clients) GetMiniPod(base basemeta.BaseMeta, podName string) (*MiniPod, error) {
    node, err := c.GetPod(base, podName)
    if err != nil {
        return nil, err
    }
    return ToMiniPod(*node), nil
}