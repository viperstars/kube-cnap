package http

import (
    "github.com/emicklei/go-restful"
    "github.com/viperstars/kube-cnap/pkg/apis/app"
)

func (s *Server) cloneApp(request *restful.Request, response *restful.Response) {
    ci := new(app.CloneInfo)
    err := request.ReadEntity(&ci)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    errs := s.dbService.CloneApp(ci)
    if len(errs) > 0 {
        s.returnResponseWithData(400, errs, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}
