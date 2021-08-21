package http

import (
    "errors"
    "github.com/emicklei/go-restful"
    app2 "github.com/viperstars/kube-cnap/pkg/apis/app"
    "github.com/viperstars/kube-cnap/pkg/apis/basemeta"
)

func (s *Server) GetAllRegions(request *restful.Request, response *restful.Response) {
    regions, err := s.dbService.GetAllRegionsAndEnvs(false)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, regions, response)
}

func (s *Server) AddOrUpdateApp(request *restful.Request, response *restful.Response) {
    intID := 0
    appID := request.PathParameter("appID")
    app := new(app2.App)
    err := request.ReadEntity(&app)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    if appID != "" {
        intID, err = s.ConvertStringToInt(appID)
        if err != nil {
            s.returnError(400, err, response)
            return
        } else {
            app.ID = intID
        }
    }
    if appID == "" {
        _app, err := s.dbService.GetAppByName(app.Name)
        if _app != nil {
            s.returnError(400, errors.New("app already exists"), response)
            return
        }
        _, err = s.dbService.CreateOrUpdateRecord(app, new(app2.App), intID)
        if err != nil {
            s.returnError(400, err, response)
            return
        }
    } else {
        err := s.dbService.UpdateApp(app)
        if err != nil {
            s.returnError(400, err, response)
            return
        }
    }
    if appID == "" {
        meta := basemeta.BaseMeta{
            App:    app.Name,
            Region: app.Region,
            Env:    app.Env,
        }
        err := s.initAppBase(meta)
        if err > 0 {
            s.returnError(400, errors.New("init app error"), response)
            return
        }
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) GetAllApps(request *restful.Request, response *restful.Response) {
    apps, err := s.dbService.GetAllApps()
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, apps, response)
}

func (s *Server) GetApp(request *restful.Request, response *restful.Response) {
    intID := 0
    appID := request.PathParameter("appID")
    app := new(app2.App)
    err := request.ReadEntity(&app)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    if appID != "" {
        intID, err = s.ConvertStringToInt(appID)
        if err != nil {
            s.returnError(400, err, response)
            return
        }
    }
    apps, err := s.dbService.GetApp(intID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithData(200, apps, response)
}

func (s *Server) GetAppByName(request *restful.Request, response *restful.Response) {
    name := request.PathParameter("name")
    if name != "" {
        app, err := s.dbService.GetAppByName(name)
        if err != nil {
            s.returnError(400, err, response)
            return
        }
        s.returnResponseWithData(200, app, response)
    } else {
        s.returnError(400, errors.New("app name is null"), response)
    }
}

func (s *Server) AddOrUpdateRegion(request *restful.Request, response *restful.Response) {
    intID := 0
    regionID := request.PathParameter("regionID")
    region := new(app2.Region)
    err := request.ReadEntity(&region)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    if regionID != "" {
        intID, err = s.ConvertStringToInt(regionID)
        if err != nil {
            s.returnError(400, err, response)
            return
        } else {
            region.ID = intID
        }
    }
    _, err = s.dbService.CreateOrUpdateRecord(region, new(app2.Region), intID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) AddOrUpdateEnvForRegion(request *restful.Request, response *restful.Response) {
    intID := 0
    envIntID := 0
    regionID := request.PathParameter("regionID")
    envID := request.PathParameter("envID")
    env := new(app2.Env)
    err := request.ReadEntity(&env)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    if regionID != "" {
        intID, err = s.ConvertStringToInt(regionID)
        if err != nil {
            s.returnError(400, err, response)
            return
        } else {
            env.RegionID = intID
        }
    }
    if envID != "" {
        envIntID, err = s.ConvertStringToInt(envID)
        if err != nil {
            s.returnError(400, err, response)
            return
        }
    }
    _, err = s.dbService.CreateOrUpdateRecord(env, new(app2.Env), envIntID)
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

func (s *Server) RemoveEnvForRegion(request *restful.Request, response *restful.Response) {
    var err error
    intID := 0
    envID := request.PathParameter("envID")
    if envID != "" {
        intID, err = s.ConvertStringToInt(envID)
        if err != nil {
            s.returnError(400, err, response)
            return
        }
    }
    _, err = s.dbService.DeleteContainerAttribute(intID, app2.Env{})
    if err != nil {
        s.returnError(400, err, response)
        return
    }
    s.returnResponseWithMessage(200, SUCCESSMESSAGE, response)
}

type DashboardResponse struct {
    Name     string      `json:"name"`
    Children interface{} `json:"children"`
    Type     string      `json:"type"`
}
