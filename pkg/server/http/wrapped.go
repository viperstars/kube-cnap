package http

import (
    "github.com/emicklei/go-restful"
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
)

func (s *Server) WrappedGetDeployment(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    deploymentName := request.PathParameter("deploymentName")
    deploymentList, err := s.k8sClients.GetMiniDeployment(base, deploymentName)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, deploymentList, response)
}

func (s *Server) WrappedGetDeployments(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    deploymentList, err := s.k8sClients.GetMiniDeploymentList(base)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, deploymentList, response)
}

func (s *Server) WrappedGetService(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    serviceName := request.PathParameter("serviceName")
    service, err := s.k8sClients.GetMiniService(base, serviceName)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, service, response)
}

func (s *Server) WrappedGetServices(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    serviceList, err := s.k8sClients.GetMiniServiceList(base)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, serviceList, response)
}

func (s *Server) WrappedGetMiniEventList(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    objectName := request.PathParameter("objectName")
    objectKind := request.PathParameter("objectKind")
    serviceList, err := s.k8sClients.GetMiniEventList(base, objectKind, objectName)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, serviceList, response)
}

func (s *Server) WrappedGetMiniNodeList(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    nodeList, err := s.k8sClients.GetMiniNodesList(base)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, nodeList, response)
}

func (s *Server) WrappedGetMiniNode(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    nodeName := request.PathParameter("nodeName")
    nodeList, err := s.k8sClients.GetMiniNode(base, nodeName)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, nodeList, response)
}

func (s *Server) WrappedGetMiniPod(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    podName := request.PathParameter("podName")
    pod, err := s.k8sClients.GetMiniPod(base, podName)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, pod, response)
}
