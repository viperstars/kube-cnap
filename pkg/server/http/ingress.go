package http

import (
    "encoding/json"
    "fmt"
    "github.com/emicklei/go-restful"
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    "github.com/viperstars/kube-cnap/pkg/apis/container"
    "github.com/viperstars/kube-cnap/pkg/apis/ingress"
)

func (s *Server) CreateOrUpdateIngress(request *restful.Request, response *restful.Response) {
    var err error
    c := new(ingress.IngressRule)
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    err = request.ReadEntity(&c)
    intID := 0
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    c.Meta = base
    id := request.PathParameter("id")
    if id != "" {
        intID, err = s.ConvertStringToInt(id)
        if err != nil {
            s.returnError(400, err, response)
            return
        }
    }
    _, err = s.dbService.CreateOrUpdateRecord(c, new(ingress.IngressRule), intID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    err = s.SyncIngress(base)

    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) SyncIngress(base basemeta.BaseMeta) error {
    igs, err := s.dbService.GetIngressRules(base, "")
    if err != nil {
        return nil
    }
    if len(igs) != 0 {
        rules := ingress.IngressRules{IngressRules: igs}
        k8sIngress := rules.ToK8sIngress(base)
        return s.k8sClients.CreateOrUpdateIngress(base, k8sIngress)
    } else {
        return s.k8sClients.DeleteIngress(base)
    }

}

func (s *Server) GetIngress(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    host := request.QueryParameter("host")
    igs, err := s.dbService.GetIngressRules(base, host)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, igs, response)
}

func (s *Server) GetIngressByID(request *restful.Request, response *restful.Response) {
    // base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    recordID := request.PathParameter("id")
    intID, err := s.ConvertStringToInt(recordID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    ig, err := s.dbService.GetIngressByID(intID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, ig, response)
}

func (s *Server) DeleteIngress(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    recordID := request.PathParameter("id")
    intID, err := s.ConvertStringToInt(recordID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    _, err = s.dbService.DeleteContainerAttribute(intID, ingress.IngressRule{})
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    err = s.SyncIngress(base)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, SUCCESSMESSAGE, response)
}

func (s *Server) DeleteConfiguration(request *restful.Request, response *restful.Response) {
    // base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    recordID := request.PathParameter("recordID")
    intID, err := s.ConvertStringToInt(recordID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    conf, err := s.dbService.GetContainerAttribute(intID, "configuration")
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    data, _ := json.Marshal(conf)
    cfg := new(container.Configuration)
    json.Unmarshal(data, &cfg)
    _, err = s.dbService.DeleteContainerAttribute(intID, container.Configuration{})
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    fmt.Println("base", cfg.Base)
    err = s.k8sClients.DeleteConfigMap(cfg.Base, cfg.Path)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, SUCCESSMESSAGE, response)
}
