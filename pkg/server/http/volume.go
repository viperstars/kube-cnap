package http

import (
    "errors"
    "fmt"
    "github.com/emicklei/go-restful"
    "github.com/viperstars/kube-cnap/pkg/apis/app"
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    "github.com/viperstars/kube-cnap/pkg/apis/pod"
)

func (s *Server) GetCFSVolumes(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    volumes, err := s.dbService.GetCFSVolumes(base)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, volumes, response)
}

func (s *Server) GetBOSVolumes(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    volumes, err := s.dbService.GetBOSVolumes(base)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, volumes, response)
}

func (s *Server) GetEmptyVolumes(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    volumes, err := s.dbService.GetEmptyVolumes(base)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, volumes, response)
}

func (s *Server) GetVolumeByID(request *restful.Request, response *restful.Response) {
    //base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    kind := request.PathParameter("kind")
    id := request.PathParameter("id")
    intID, err := s.ConvertStringToInt(id)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    data, err := s.dbService.GetVolumeByID(intID, kind)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, data, response)
}

func (s *Server) CreateOrUpdateBOSVolumes(request *restful.Request, response *restful.Response) {
    // base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    id := request.PathParameter("name")
    var err error
    intID := 0
    if id != "" {
        intID, err = s.ConvertStringToInt(id)
        if err != nil {
            s.returnError(400, err, response)
            return
        }
    }
    volume := &pod.BOSVolume{}
    err = request.ReadEntity(volume)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    //volume.Meta = base
    _, err = s.dbService.CreateOrUpdateRecord(volume, new(pod.BOSVolume), intID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, SUCCESSMESSAGE, response)
}

func (s *Server) CreateOrUpdateEmptyVolumes(request *restful.Request, response *restful.Response) {
    //base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    id := request.PathParameter("id")
    var err error
    intID := 0
    if id != "" {
        intID, err = s.ConvertStringToInt(id)
        if err != nil {
            s.returnError(400, err, response)
            return
        }
    }
    volume := &pod.EmptyVolume{}
    err = request.ReadEntity(volume)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    // volume.Meta = base
    _, err = s.dbService.CreateOrUpdateRecord(volume, new(pod.EmptyVolume), intID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, SUCCESSMESSAGE, response)
}

func (s *Server) CreateOrUpdateCFSVolumes(request *restful.Request, response *restful.Response) {
    //base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    id := request.PathParameter("id")
    var err error
    intID := 0
    if id != "" {
        intID, err = s.ConvertStringToInt(id)
        if err != nil {
            s.returnError(400, err, response)
            return
        }
    }
    volume := &pod.CFSVolume{}
    err = request.ReadEntity(volume)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    //volume.Meta = base
    _, err = s.dbService.CreateOrUpdateRecord(volume, new(pod.CFSVolume), intID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, SUCCESSMESSAGE, response)
}

func (s *Server) DeleteVolume(request *restful.Request, response *restful.Response) {
    // base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    id := request.PathParameter("id")
    kind := request.PathParameter("kind")
    var err error
    var table interface{}
    intID := 0
    switch kind {
    case "bos":
        table = &pod.BOSVolume{}
    case "empty":
        table = &pod.EmptyVolume{}
    case "cfs":
        table = &pod.CFSVolume{}
    default:
        s.returnError(400, errors.New("error kind, should be"), response)
        return
    }
    if id != "" {
        intID, err = s.ConvertStringToInt(id)
        if err != nil {
            s.returnError(400, err, response)
            return
        }
    }
    _, err = s.dbService.DeleteContainerAttribute(intID, table)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, SUCCESSMESSAGE, response)
}

func (s *Server) AddHistory(request *restful.Request, response *restful.Response) {
    // base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    fmt.Println("in")
    history := new(app.History)
    err := request.ReadEntity(&history)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    _, err = s.dbService.CreateOrUpdateRecord(history, app.History{}, 0)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, SUCCESSMESSAGE, response)
}

func (s *Server) GetHistories(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    his, err := s.dbService.GetHistories(base)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, his, response)
}

func (s *Server) CreateSecret2(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    err := s.k8sClients.CreateSecret2(base)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, "success", response)
}

func (s *Server) initZone(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    errs := s.initAppBase(base)
    if errs != 0 {
        err := fmt.Errorf("init region and env error")
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, "success", response)
}

func (s *Server) getNodeSelector(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    ns, err := s.dbService.GetNodeSelector(base)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, ns, response)
}

func (s *Server) updateNodeSelector(request *restful.Request, response *restful.Response) {
    id := request.PathParameter("id")
    var err error
    intID := 0
    if id != "" {
        intID, err = s.ConvertStringToInt(id)
        if err != nil {
            s.returnError(400, err, response)
            return
        }
    }
    ns := &pod.NodeSelector{}
    err = request.ReadEntity(ns)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    err = s.dbService.UpdateNodeSelector(intID, ns)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) deleteNodeSelector(request *restful.Request, response *restful.Response) {
    id := request.PathParameter("id")
    var err error
    intID := 0
    if id != "" {
        intID, err = s.ConvertStringToInt(id)
        if err != nil {
            s.returnError(400, err, response)
            return
        }
    }
    err = s.dbService.DeleteNodeSelector(intID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) insertNodeSelector(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    ns := &pod.NodeSelector{}
    err := request.ReadEntity(ns)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    err = s.dbService.InsertNodeSelector(base, ns)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) getHostAlias(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    hs, err := s.dbService.GetHostAlias(base)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, hs, response)
}

func (s *Server) updateHostAlias(request *restful.Request, response *restful.Response) {
    id := request.PathParameter("id")
    var err error
    intID := 0
    if id != "" {
        intID, err = s.ConvertStringToInt(id)
        if err != nil {
            s.returnError(400, err, response)
            return
        }
    }
    hs := &pod.HostAlias{}
    err = request.ReadEntity(hs)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    err = s.dbService.UpdateHostAlias(intID, hs)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) deleteHostAlias(request *restful.Request, response *restful.Response) {
    id := request.PathParameter("id")
    var err error
    intID := 0
    if id != "" {
        intID, err = s.ConvertStringToInt(id)
        if err != nil {
            s.returnError(400, err, response)
            return
        }
    }
    err = s.dbService.DeleteNodeSelector(intID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) insertHostAlias(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    hs := &pod.HostAlias{}
    err := request.ReadEntity(hs)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    err = s.dbService.InsertHostAlias(base, hs)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}


