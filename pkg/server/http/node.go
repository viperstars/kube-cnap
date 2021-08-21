package http

import (
    "github.com/emicklei/go-restful"
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
)

func (s *Server) GetNodesHandler(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    nodes, err := s.k8sClients.GetNodes(base)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, nodes, response)
}

func (s *Server) UpdateNodeLabels(request *restful.Request, response *restful.Response) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    nodeName := request.PathParameter("nodeName")
    lbs := make(map[string]string)
    err := request.ReadEntity(&lbs)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    err = s.k8sClients.UpdateLabelsToNode(base, nodeName, lbs)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, SUCCESSMESSAGE, response)
}
