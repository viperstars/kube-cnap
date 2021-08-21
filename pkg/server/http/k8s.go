package http

import (
    "errors"
    "github.com/emicklei/go-restful"
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    "github.com/viperstars/kube-cnap/pkg/apis/container"
    "github.com/viperstars/kube-cnap/pkg/apis/requests"
    "github.com/viperstars/kube-cnap/pkg/apis/service"
    "github.com/viperstars/kube-cnap/pkg/client"
    v1 "k8s.io/api/core/v1"
    "k8s.io/client-go/tools/remotecommand"
)

var BADMETAINFO = errors.New("bad meta info, url prefix should be /app/region/env")

func (s *Server) CreateDeploymentHandler(request *restful.Request, response *restful.Response) {
    c := new(requests.CreateDeploymentRequest)
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    err := request.ReadEntity(&c)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    err = s.CreateDeployment(base, c)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    cg, err := s.dbService.GetContainerGroupBase(base)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    cg.MainBase.Image = c.Image
    _, err = s.dbService.CreateOrUpdateRecord(cg.MainBase, new(container.Base), cg.MainBase.ID)
    if err != nil  {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) GetDeploymentHandler(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    deploymentName := request.PathParameter("deploymentName")
    deployment, err := s.k8sClients.GetDeployment(base, deploymentName)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, deployment, response)
}

func (s *Server) GetDeploymentsHandler(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    deploymentList, err := s.k8sClients.GetDeployments(base, base.GetLabels())
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, deploymentList, response)
}

func (s *Server) GetServiceHandler(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    serviceName := request.PathParameter("serviceName")
    svc, err := s.k8sClients.GetService(base, serviceName)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, svc, response)
}

func (s *Server) GetServicesHandler(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    serviceList, err := s.k8sClients.GetServices(base, base.GetLabels())
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, serviceList, response)
}

func (s *Server) CreateServiceHandler(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    c := new(service.CreateServiceRequest)
    err := request.ReadEntity(&c)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    err = s.CreateService(base, c)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, SUCCESSMESSAGE, response)
}

func (s *Server) UpdateDeploymentHandler(request *restful.Request, response *restful.Response) {
    updateRequest := new(requests.UpdateDeploymentRequest)
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    deploymentName := request.PathParameter("deploymentName")
    err := request.ReadEntity(&updateRequest)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    if updateRequest.VersionOnly {
        err = s.k8sClients.UpdateDeployment(base, deploymentName, updateRequest.Image)
        if err != nil {
            s.returnError(400, err, response)
            return
        }
    } else {
        err = s.RolloutDeployment(base, deploymentName, updateRequest.Image)
        if err != nil {
            s.returnError(400, err, response)
            return
        }
    }
    cg, err := s.dbService.GetContainerGroupBase(base)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    cg.MainBase.Image = updateRequest.Image
    _, err = s.dbService.CreateOrUpdateRecord(cg.MainBase, new(container.Base), cg.MainBase.ID)
    if err != nil  {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, SUCCESSMESSAGE, response)
}

func (s *Server) RolloutDeploymentHandler(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    deploymentName := request.PathParameter("deploymentName")
    dep, err := s.k8sClients.GetDeployment(base, deploymentName)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    if len(dep.Spec.Template.Spec.Containers) > 0 {
        err := s.RolloutDeployment(base, deploymentName, dep.Spec.Template.Spec.Containers[0].Image)
        if err != nil {
            s.returnError(400, err, response)
            return
        }
        s.returnResponseWithData(200, SUCCESSMESSAGE, response)
    } else {
        s.returnError(400, errors.New("dep.Spec.Template.Spec.Containers <= 0"), response)
        return
    }
}

func (s *Server) ScaleDeploymentHandler(request *restful.Request, response *restful.Response) {
    updateRequest := new(requests.ScaleDeploymentRequest)
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    podName := request.PathParameter("deploymentName")
    err := request.ReadEntity(&updateRequest)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    err = s.k8sClients.ScaleDeployment(base, podName, updateRequest.Replicas)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, SUCCESSMESSAGE, response)
}

func (s *Server) DeleteDeploymentHandler(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    deploymentName := request.PathParameter("deploymentName")
    err := s.k8sClients.DeleteDeployment(base, deploymentName)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, SUCCESSMESSAGE, response)
}

func (s *Server) DeleteServiceHandler(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    serviceName := request.PathParameter("serviceName")
    err := s.k8sClients.DeleteService(base, serviceName)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, SUCCESSMESSAGE, response)
}

func (s *Server) DeletePodHandler(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    podName := request.PathParameter("podName")
    err := s.k8sClients.DeletePod(base, podName)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, SUCCESSMESSAGE, response)
}

func (s *Server) GetEventsHandler(request *restful.Request, response *restful.Response) {
    var eventList *v1.EventList
    var err error
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    objectName := request.PathParameter("objectName")
    objectKind := request.PathParameter("objectKind")
    if objectKind == "node" {
        eventList, err = s.k8sClients.GetEventsForNode(base, objectName, false)
    } else {
        eventList, err = s.k8sClients.GetEventsForObject(base, objectName)
    }
    if err != nil {
        s.returnError(500, err, response)
        return
    }
    s.returnResponseWithData(200, eventList, response)
}

func (s *Server) handleExecShell(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    sessionId, err := genTerminalSessionId()
    if err != nil {
        s.returnError(500, err, response)
        return
    }
    c, config, _ := s.k8sClients.GetClientAndConfig(base)
    terminalSessions.Set(sessionId, TerminalSession{
        id:       sessionId,
        bound:    make(chan error),
        sizeChan: make(chan remotecommand.TerminalSize),
    })
    go WaitForTerminal(c, config, base.App, request, sessionId)
    s.returnResponseWithData(200, TerminalResponse{Id: sessionId}, response)
    //response.WriteAsJson(TerminalResponse{Id: sessionId})
}

func (s *Server) dashboard(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    deps, err := s.getDeploymentsAndDetails(base)
    if err != nil {
        s.returnError(500, err, response)
        return
    }
    allDep := make([]*client.MiniDeployment, 0)
    for _, d := range deps {
        allDep = append(allDep, d)
    }
    data := &DashboardResponse{
        Name:     "外部服务",
        Type:     "外部服务",
    }
    services, err := s.k8sClients.GetMiniServiceList(base)
    if len(services) > 0 {
        for i:=0; i<len(services); i+=1 {
            services[i].IP = services[i].ClusterIP
            if services[i].DeploymentName == "all" {
                services[i].Children = allDep
            } else if v, ok := deps[services[i].DeploymentName]; ok {
                services[i].Children = []interface{}{v}
            }
        }
        data.Children = services
    }
    s.returnResponseWithData(200, data, response)
}

func (s *Server) getDeploymentsAndDetails(base basemeta.BaseMeta) (map[string]*client.MiniDeployment, error) {
    miniDeps := make([]*client.MiniDeployment, 0)
    depsMap := make(map[string]*client.MiniDeployment)
    deps, err := s.k8sClients.GetDeployments(base, nil)
    if err != nil {
        return depsMap, err
    }
    for _, d := range deps.Items {
        miniDep, err := s.k8sClients.GetMiniDeployment(base, d.Name)
        if err == nil {
            miniDep.Children = miniDep.NewestRs.PodList
            miniDeps = append(miniDeps, miniDep)
        }
    }
    for _, mini := range miniDeps {
        depsMap[mini.Name] = mini
    }
    return depsMap, nil
}
