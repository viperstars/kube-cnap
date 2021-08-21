package http

import (
    "errors"
    "fmt"
    "github.com/emicklei/go-restful"
    app2 "github.com/viperstars/kube-cnap/pkg/apis/app"
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
    "strings"
)

func (s *Server) getMeta(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
    app := request.PathParameter("app")
    region := request.PathParameter("region")
    env := request.PathParameter("env")
    if app != "" && region != "" && env != "" {
        meta := basemeta.BaseMeta{
            App:    app,
            Region: region,
            Env:    env,
        }
        request.SetAttribute("meta", meta)
        chain.ProcessFilter(request, response)
    } else {
        s.returnError(400, errors.New("bad url scheme, need to include app, region, env"), response)
        return
    }
}

func (s *Server) currentApp(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    app, err := s.dbService.GetAppByName(base.App)
    if err != nil {
        s.returnError(400, errors.New("bad url scheme, need to include app, region, env"), response)
        return
    } else {
        request.SetAttribute("app", app)
        chain.ProcessFilter(request, response)
    }
}

func (s *Server) validateRoles(roles []string) func(request *restful.Request, response *restful.Response,
    chain *restful.FilterChain) {
    return func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
        app, _ := request.Attribute("app").(app2.App)
        user := request.HeaderParameter("user")
        if user == "" {
            s.returnError(400, errors.New("user not in headers"), response)
            return
        } else if s.stringInSlice(user, app.ReturnUsers(roles)) {
            chain.ProcessFilter(request, response)
        } else {
            msg := fmt.Sprintf("user: %s has no permission", user)
            s.returnError(403, errors.New(msg), response)
        }
    }
}


func (s *Server) validateMeta(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
    meta, ok := request.Attribute("meta").(basemeta.BaseMeta)
    if ok {
        err := s.k8sClients.EnsureGet(meta)
        if err != nil {
            s.returnError(400, errors.New("bad url scheme, need to include app, region, env"), response)
            return
        }
        chain.ProcessFilter(request, response)
    } else {
        s.returnError(400, errors.New("bad url scheme, need to include app, region, env"), response)
        return
    }
}

func (s *Server) validatePathParamPrefix(param string) func(request *restful.Request, response *restful.Response,
    chain *restful.FilterChain) () {
    return func(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
        base, _ := request.Attribute("meta").(basemeta.BaseMeta)
        prefix := base.GetPrefix()
        p := request.PathParameter(param)
        if p != "" && strings.HasPrefix(p, prefix) {
            chain.ProcessFilter(request, response)
        } else {
            s.returnError(400, errors.New("bad param, param: "+p), response)
            return
        }
    }
}

func (s *Server) validateCapacity(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
    base, _ := request.Attribute("meta").(basemeta.BaseMeta)
    _, err := s.k8sClients.GetResourceQuota(base)
    if err != nil {
        s.returnError(400, errors.New("capacity not set"), response)
        return
    } else {
        chain.ProcessFilter(request, response)
    }
}