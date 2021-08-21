package http

import (
    "github.com/emicklei/go-restful"
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    "github.com/viperstars/kube-cnap/pkg/apis/container"
    "github.com/viperstars/kube-cnap/pkg/apis/resourcequota"
)

type IDs []int

func (s *Server) GetContainerBase(request *restful.Request, response *restful.Response) {
    var intID int
    var err error
    id := request.PathParameter("id")
    if id != "" {
        intID, err = s.ConvertStringToInt(id)
        if err != nil {
            s.returnError(400, err, response)
            return
        }
    }
    containerBase, err := s.dbService.GetContainerBase(intID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, containerBase, response)
}

func (s *Server) AddContainerBase(request *restful.Request, response *restful.Response) {
    var err error
    c := new(container.Base)
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    kind := request.PathParameter("kind")
    err = request.ReadEntity(&c)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    err = s.dbService.AddContainer(base, kind, c)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) UpdateContainerBase(request *restful.Request, response *restful.Response) {
    var err error
    c := new(container.Base)
    id := request.PathParameter("id")
    intID, err := s.ConvertStringToInt(id)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    err = request.ReadEntity(&c)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    c.ID = intID
    _, err = s.dbService.CreateOrUpdateRecord(c, new(container.Base), intID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) UpdateInitOrSidecarContainers(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    kind := request.PathParameter("kind")
    ids := new(IDs)
    err := request.ReadEntity(&ids)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    err = s.dbService.UpdateInitOrSidecarContainer(base, kind, *ids)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) RemoveMainContainer(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    err := s.dbService.UpdateMainContainer(base, 0)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) GetContainerGroup(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    cg, err := s.dbService.GetContainerGroupBase(base)
    if err != nil {
        s.returnError(404, err, response)
        return
    }
    s.returnResponseWithData(200, cg, response)
}

func (s *Server) RemoveContainerBase(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    kind := request.PathParameter("kind")
    id := request.PathParameter("id")
    intID, err := s.ConvertStringToInt(id)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    if kind == "main" {
        err = s.dbService.RemoveMainContainer(base)
    } else {
        err = s.dbService.RemoveInitOrSidecarContainerByID(base, kind, intID)
    }
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) GetResourceQuota(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    rq, err := s.k8sClients.GetResourceQuota(base)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, rq, response)
}

func (s *Server) WrappedGetResourceQuota(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    rq, err := s.k8sClients.GetResourceQuota(base)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    crq := s.ToCapacitys(base, rq)
    s.returnResponseWithData(200, crq, response)
}

func (s *Server) UpdateResourceQuota(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    rq := new(resourcequota.ResourceQuota)
    err := request.ReadEntity(rq)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    rq.Base = base
    k8sRq := rq.ToK8sResourceQuota()
    _, err = s.k8sClients.UpdateResourceQuota(base, k8sRq)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) CreateResourceQuota(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    rq := new(resourcequota.ResourceQuota)
    err := request.ReadEntity(rq)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    rq.Base = base
    k8sRq := rq.ToK8sResourceQuota()
    err = s.k8sClients.AddResourceQuota(base, k8sRq)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}
